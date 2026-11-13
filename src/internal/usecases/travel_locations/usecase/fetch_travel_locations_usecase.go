package usecase

import (
	"context"

	api "platform-go/generated/api"
	"platform-go/internal/infrastructure/qsdto"
	travellocationqueryservice "platform-go/internal/usecases/travel_locations/queryservice"
)

type FetchTravelLocationsUsecase struct {
	queryService travellocationqueryservice.TravelLocationQueryService
}

func NewFetchTravelLocationsUsecase(queryService travellocationqueryservice.TravelLocationQueryService) *FetchTravelLocationsUsecase {
	return &FetchTravelLocationsUsecase{queryService: queryService}
}

func (uc *FetchTravelLocationsUsecase) Execute(ctx context.Context, params api.FetchTravelLocationsParams) (*api.FetchTravelLocationsResponse, error) {
	var status *string
	if params.Status != nil {
		s := string(*params.Status)
		status = &s
	}

	filter := qsdto.TravelLocationFilter{
		Status:      status,
		CountryCode: params.CountryCode,
	}

	locations, total, err := uc.queryService.FetchTravelLocations(ctx, filter)
	if err != nil {
		return nil, err
	}

	items := make([]api.FetchTravelLocationsResponseDataInner, len(locations))
	for i, loc := range locations {
		items[i] = api.FetchTravelLocationsResponseDataInner{
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
		}
	}

	return &api.FetchTravelLocationsResponse{
		Data:  items,
		Total: total,
	}, nil
}
