package travel_memory_controller_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	api "platform-go/generated/api"
	"platform-go/internal/infrastructure/qsdto"
)

func TestFetchTravelMemoryByID_Success(t *testing.T) {
	desc := "Kỷ niệm tuyệt vời tại Đà Lạt"
	dayDesc := "Ngày đầu tiên đặt chân đến Đà Lạt"
	startDate := time.Date(2025, 3, 10, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 3, 13, 0, 0, 0, 0, time.UTC)

	qs := &travelMemoryQueryService{
		memory: &qsdto.TravelMemoryDto{
			ID:             1,
			LocationID:     1,
			LocationName:   "Đà Lạt",
			CountryCode:    "VN",
			CountryName:    "Vietnam",
			Title:          "Mùa hoa dã quỳ trên cao nguyên Đà Lạt",
			StartDate:      startDate,
			EndDate:        endDate,
			Description:    &desc,
			CoverImagePath: "travel_memories/a1b2c3d4-e5f6-7890-abcd-ef1234567890.jpg",
			Days: []qsdto.TravelMemoryDayDto{
				{
					ID:          1,
					MemoryID:    1,
					Date:        startDate,
					Description: &dayDesc,
					Photos: []qsdto.TravelMemoryPhotoDto{
						{
							ID:          1,
							MemoryDayID: 1,
							ImagePath:   "travel_memories/f1e2d3c4-b5a6-7890-1234-567890abcdef.jpg",
							SortOrder:   1,
						},
					},
				},
			},
		},
	}

	router := newTravelMemoryRouter(t, qs)
	request := httptest.NewRequest(http.MethodGet, "/travel-memories/1", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", response.Code, response.Body.String())
	}

	var res api.FetchTravelMemoryByIDResponse
	if err := json.Unmarshal(response.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	}

	if res.Data.Id != 1 || res.Data.Title != "Mùa hoa dã quỳ trên cao nguyên Đà Lạt" {
		t.Fatalf("unexpected travel memory data: %#v", res.Data)
	}
	if res.Data.Location.Id != 1 || res.Data.Location.LocationName != "Đà Lạt" || res.Data.Location.CountryCode != "VN" || res.Data.Location.CountryName != "Vietnam" {
		t.Fatalf("unexpected location data: %#v", res.Data.Location)
	}
	if len(res.Data.Days) != 1 || len(res.Data.Days[0].Photos) != 1 {
		t.Fatalf("unexpected days or photos: %#v", res.Data.Days)
	}
	if res.Data.Days[0].Photos[0].ImagePath != "travel_memories/f1e2d3c4-b5a6-7890-1234-567890abcdef.jpg" {
		t.Fatalf("unexpected photo image path: %s", res.Data.Days[0].Photos[0].ImagePath)
	}

	rawJSON := response.Body.String()
	for _, camelField := range []string{
		`"locationId"`,
		`"startDate"`,
		`"endDate"`,
		`"coverImagePath"`,
		`"location"`,
		`"locationName"`,
		`"countryCode"`,
		`"countryName"`,
		`"days"`,
		`"memoryId"`,
		`"photos"`,
		`"memoryDayId"`,
		`"imagePath"`,
		`"sortOrder"`,
	} {
		if !strings.Contains(rawJSON, camelField) {
			t.Fatalf("expected JSON response to contain camelCase field %s, got: %s", camelField, rawJSON)
		}
	}
}

func TestFetchTravelMemoryByID_NotFound(t *testing.T) {
	qs := &travelMemoryQueryService{memory: nil}
	router := newTravelMemoryRouter(t, qs)
	request := httptest.NewRequest(http.MethodGet, "/travel-memories/999", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d, body: %s", response.Code, response.Body.String())
	}
}

func TestFetchTravelMemoryByID_InvalidID(t *testing.T) {
	qs := &travelMemoryQueryService{}
	router := newTravelMemoryRouter(t, qs)
	request := httptest.NewRequest(http.MethodGet, "/travel-memories/not-a-number", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body: %s", response.Code, response.Body.String())
	}
}

func TestFetchTravelMemoryByID_DatabaseError(t *testing.T) {
	qs := &travelMemoryQueryService{err: errors.New("db query failed")}
	router := newTravelMemoryRouter(t, qs)
	request := httptest.NewRequest(http.MethodGet, "/travel-memories/1", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d, body: %s", response.Code, response.Body.String())
	}
}

