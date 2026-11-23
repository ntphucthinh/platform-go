package usecase

import (
	"context"
	"errors"

	openapi_types "github.com/oapi-codegen/runtime/types"

	api "platform-go/generated/api"
	travellocationqueryservice "platform-go/internal/usecases/travel_locations/queryservice"
)

var ErrTravelLocationNotFound = errors.New("travel location not found")

type FetchTravelLocationByIDUsecase struct {
	queryService travellocationqueryservice.TravelLocationQueryService
}

func NewFetchTravelLocationByIDUsecase(queryService travellocationqueryservice.TravelLocationQueryService) *FetchTravelLocationByIDUsecase {
	return &FetchTravelLocationByIDUsecase{queryService: queryService}
}

func (uc *FetchTravelLocationByIDUsecase) Execute(ctx context.Context, id int64) (*api.FetchTravelLocationByIDResponse, error) {
	loc, err := uc.queryService.FetchTravelLocationByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if loc == nil {
		return nil, ErrTravelLocationNotFound
	}

	memories := make([]api.FetchTravelLocationByIDResponseDataMemoriesInner, len(loc.Memories))
	for i, m := range loc.Memories {
		memories[i] = api.FetchTravelLocationByIDResponseDataMemoriesInner{
			Id:             m.ID,
			Title:          m.Title,
			CoverImagePath: m.CoverImagePath,
			StartDate:      openapi_types.Date{Time: m.StartDate},
			EndDate:        openapi_types.Date{Time: m.EndDate},
			PhotoCount:     m.PhotoCount,
		}
	}

	itineraries := make([]api.FetchTravelLocationByIDResponseDataItinerariesInner, len(loc.Itineraries))
	for i, it := range loc.Itineraries {
		itineraries[i] = api.FetchTravelLocationByIDResponseDataItinerariesInner{
			Id:       it.ID,
			Title:    it.Title,
			Duration: it.Duration,
		}
	}

	return &api.FetchTravelLocationByIDResponse{
		Data: api.FetchTravelLocationByIDResponseData{
			Id:             loc.ID,
			CountryCode:    loc.CountryCode,
			CountryName:    loc.CountryName,
			LocationName:   loc.LocationName,
			Latitude:       loc.Latitude,
			Longitude:      loc.Longitude,
			ImagePath:      loc.ImagePath,
			Status:         api.TravelLocationStatus(loc.Status),
			Description:    loc.Description,
			Notes:          loc.Notes,
			GoogleMapsUrl:  loc.GoogleMapsURL,
			VisitCount:     loc.VisitCount,
			ItineraryCount: loc.ItineraryCount,
			Memories:       memories,
			Itineraries:    itineraries,
		},
	}, nil
}
