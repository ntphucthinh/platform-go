package usecase

import (
	"context"
	"errors"

	api "platform-go/generated/api"
	travellocationqueryservice "platform-go/internal/usecases/travel_locations/queryservice"
)

var ErrTravelLocationNotFound = errors.New("travel location not found")

type FetchTravelLocationByIDUseCase struct {
	queryService travellocationqueryservice.TravelLocationQueryService
}

func NewFetchTravelLocationByIDUseCase(queryService travellocationqueryservice.TravelLocationQueryService) *FetchTravelLocationByIDUseCase {
	return &FetchTravelLocationByIDUseCase{queryService: queryService}
}

func (uc *FetchTravelLocationByIDUseCase) Execute(ctx context.Context, id int64) (*api.FetchTravelLocationByIDResponse, error) {
	loc, err := uc.queryService.FetchTravelLocationByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if loc == nil {
		return nil, ErrTravelLocationNotFound
	}

	return &api.FetchTravelLocationByIDResponse{
		Data: api.FetchTravelLocationsResponseDataInner{
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
		},
	}, nil
}

