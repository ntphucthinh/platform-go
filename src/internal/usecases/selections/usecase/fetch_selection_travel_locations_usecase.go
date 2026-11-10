package usecase

import (
	"context"

	api "platform-go/generated/api"
	selectionqueryservice "platform-go/internal/usecases/selections/queryservice"
)

type FetchSelectionTravelLocationsUseCase struct {
	queryService selectionqueryservice.SelectionQueryService
}

func NewFetchSelectionTravelLocationsUseCase(queryService selectionqueryservice.SelectionQueryService) *FetchSelectionTravelLocationsUseCase {
	return &FetchSelectionTravelLocationsUseCase{
		queryService: queryService,
	}
}

func (uc *FetchSelectionTravelLocationsUseCase) Execute(ctx context.Context, countryCode *string) (*api.SelectionTravelLocationResponse, error) {
	items, err := uc.queryService.FetchTravelLocations(ctx, countryCode)
	if err != nil {
		return nil, err
	}

	data := make([]api.SelectionTravelLocationResponseDataInner, len(items))
	for i, item := range items {
		data[i] = api.SelectionTravelLocationResponseDataInner{
			Id:   item.ID,
			Name: item.Name,
		}
	}

	return &api.SelectionTravelLocationResponse{
		Data: data,
	}, nil
}

