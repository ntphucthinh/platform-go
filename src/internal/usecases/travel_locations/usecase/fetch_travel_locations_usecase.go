package usecase

import (
	"context"

	api "platform-go/generated/api"
	"platform-go/internal/infrastructure/qsdto"
	travellocationqueryservice "platform-go/internal/usecases/travel_locations/queryservice"
)

type FetchTravelLocationsUseCase struct {
	queryService travellocationqueryservice.TravelLocationQueryService
}

func NewFetchTravelLocationsUseCase(queryService travellocationqueryservice.TravelLocationQueryService) *FetchTravelLocationsUseCase {
	return &FetchTravelLocationsUseCase{queryService: queryService}
}

func (uc *FetchTravelLocationsUseCase) Execute(ctx context.Context, params api.FetchTravelLocationsParams) (*api.FetchTravelLocationsResponse, error) {
	var status *string
	if params.Status != nil {
		s := string(*params.Status)
		status = &s
	}

	filter := qsdto.TravelLocationFilter{
		Status:      status,
		Name:        params.Name,
		StateCode:   params.StateCode,
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
			StateCode:      loc.StateCode,
			Name:           loc.Name,
			Address:        loc.Address,
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

