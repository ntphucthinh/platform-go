package travel_memory_controller_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	api "platform-go/generated/api"
	"platform-go/internal/controller_test/shared"
	"platform-go/internal/infrastructure/dbutil"
	"platform-go/internal/infrastructure/qsdto"
	"platform-go/internal/interfaces/controller"
	travelmemoryqueryservice "platform-go/internal/usecases/travel_memories/queryservice"
	travelmemoryusecase "platform-go/internal/usecases/travel_memories/usecase"
)

type travelMemoryQueryService struct {
	memories       []qsdto.TravelMemoryDto
	memory         *qsdto.TravelMemoryDto
	lastPerPage    int
	lastCountryCode *string
	lastLocationID *int64
	err            error
}

func (s *travelMemoryQueryService) FetchTravelMemories(ctx context.Context, page, perPage int, countryCode *string, locationID *int64) (*dbutil.Paginator[qsdto.TravelMemoryDto], error) {
	s.lastPerPage = perPage
	s.lastCountryCode = countryCode
	s.lastLocationID = locationID
	if s.err != nil {
		return nil, s.err
	}
	return &dbutil.Paginator[qsdto.TravelMemoryDto]{
		Data: s.memories,
		PaginatorInfo: api.PaginatorInfo{
			Count:        len(s.memories),
			CurrentPage:  page,
			HasMorePages: false,
			LastPage:     1,
			PerPage:      perPage,
			Total:        int64(len(s.memories)),
		},
	}, nil
}

func (s *travelMemoryQueryService) FetchTravelMemoryByID(ctx context.Context, id int64) (*qsdto.TravelMemoryDto, error) {
	return s.memory, s.err
}

var _ travelmemoryqueryservice.TravelMemoryQueryService = (*travelMemoryQueryService)(nil)

func newTravelMemoryRouter(t *testing.T, queryService travelmemoryqueryservice.TravelMemoryQueryService) http.Handler {
	t.Helper()
	fetchUsecase := travelmemoryusecase.NewFetchTravelMemoriesUseCase(queryService)
	fetchByIDUsecase := travelmemoryusecase.NewFetchTravelMemoryByIDUseCase(queryService)
	handler := controller.NewTravelMemoryController(fetchUsecase, fetchByIDUsecase)
	serverController := controller.NewServerController(nil, nil, nil, nil, handler, nil)
	router := shared.NewSetupRouter(t)
	api.RegisterHandlers(router, serverController)
	return router
}

func TestFetchTravelMemories(t *testing.T) {
	desc := "Kỷ niệm tuyệt vời tại Đà Lạt"
	dayDesc := "Ngày đầu tiên đặt chân đến Đà Lạt"
	startDate := time.Date(2025, 3, 10, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 3, 13, 0, 0, 0, 0, time.UTC)

	qs := &travelMemoryQueryService{
		memories: []qsdto.TravelMemoryDto{
			{
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
		},
	}

	router := newTravelMemoryRouter(t, qs)
	request := httptest.NewRequest(http.MethodGet, "/travel-memories?page=1", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	if qs.lastPerPage != 10 {
		t.Fatalf("expected fixed perPage = 10, got %d", qs.lastPerPage)
	}

	var body api.FetchTravelMemoriesResponse
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if len(body.Data) != 1 || len(body.Data[0].Days) != 1 || len(body.Data[0].Days[0].Photos) != 1 {
		t.Fatalf("response = %#v", body)
	}
	if got := body.Data[0].CoverImagePath; got != "travel_memories/a1b2c3d4-e5f6-7890-abcd-ef1234567890.jpg" {
		t.Fatalf("cover image path = %q", got)
	}
	if got := body.Data[0].Days[0].Photos[0].ImagePath; got != "travel_memories/f1e2d3c4-b5a6-7890-1234-567890abcdef.jpg" {
		t.Fatalf("photo image path = %q", got)
	}
	if body.PaginatorInfo.CurrentPage != 1 || body.PaginatorInfo.PerPage != 10 || body.PaginatorInfo.Count != 1 || body.PaginatorInfo.Total != 1 {
		t.Fatalf("paginatorInfo = %#v", body.PaginatorInfo)
	}

	rawJSON := response.Body.String()
	for _, camelField := range []string{
		`"locationId"`,
		`"locationName"`,
		`"countryCode"`,
		`"startDate"`,
		`"endDate"`,
		`"coverImagePath"`,
		`"days"`,
		`"memoryId"`,
		`"photos"`,
		`"memoryDayId"`,
		`"imagePath"`,
		`"sortOrder"`,
		`"paginatorInfo"`,
	} {
		if !strings.Contains(rawJSON, camelField) {
			t.Fatalf("expected JSON response to contain camelCase field %s, got: %s", camelField, rawJSON)
		}
	}
}

func TestFetchTravelMemories_FilterByLocationID(t *testing.T) {
	qs := &travelMemoryQueryService{}
	router := newTravelMemoryRouter(t, qs)

	request := httptest.NewRequest(http.MethodGet, "/travel-memories?page=1&locationId=5", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	if qs.lastLocationID == nil || *qs.lastLocationID != 5 {
		t.Fatalf("expected lastLocationID = 5, got: %v", qs.lastLocationID)
	}
}

func TestFetchTravelMemories_FilterByCountryCode(t *testing.T) {
	qs := &travelMemoryQueryService{}
	router := newTravelMemoryRouter(t, qs)

	request := httptest.NewRequest(http.MethodGet, "/travel-memories?page=1&countryCode=JP", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	if qs.lastCountryCode == nil || *qs.lastCountryCode != "JP" {
		t.Fatalf("expected lastCountryCode = JP, got: %v", qs.lastCountryCode)
	}
}

func TestFetchTravelMemories_FilterByCountryCodeSnakeCase(t *testing.T) {
	qs := &travelMemoryQueryService{}
	router := newTravelMemoryRouter(t, qs)

	request := httptest.NewRequest(http.MethodGet, "/travel-memories?page=1&country_code=vn", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	if qs.lastCountryCode == nil || *qs.lastCountryCode != "vn" {
		t.Fatalf("expected lastCountryCode = vn, got: %v", qs.lastCountryCode)
	}
}

func TestFetchTravelMemories_FilterByBothCountryAndLocation(t *testing.T) {
	qs := &travelMemoryQueryService{}
	router := newTravelMemoryRouter(t, qs)

	request := httptest.NewRequest(http.MethodGet, "/travel-memories?page=1&countryCode=VN&locationId=10", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	if qs.lastCountryCode == nil || *qs.lastCountryCode != "VN" {
		t.Fatalf("expected lastCountryCode = VN, got: %v", qs.lastCountryCode)
	}
	if qs.lastLocationID == nil || *qs.lastLocationID != 10 {
		t.Fatalf("expected lastLocationID = 10, got: %v", qs.lastLocationID)
	}
}

func TestFetchTravelMemories_MissingRequiredParams(t *testing.T) {
	router := newTravelMemoryRouter(t, &travelMemoryQueryService{})
	request := httptest.NewRequest(http.MethodGet, "/travel-memories", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d, body = %s", response.Code, response.Body.String())
	}
}

func TestFetchTravelMemories_DatabaseError(t *testing.T) {
	router := newTravelMemoryRouter(t, &travelMemoryQueryService{err: errors.New("database unavailable")})
	request := httptest.NewRequest(http.MethodGet, "/travel-memories?page=1", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d, body = %s", response.Code, response.Body.String())
	}
}

