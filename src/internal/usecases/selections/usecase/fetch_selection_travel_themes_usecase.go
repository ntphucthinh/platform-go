package usecase

import (
	"context"

	api "platform-go/generated/api"
	selectionqueryservice "platform-go/internal/usecases/selections/queryservice"
)

type FetchSelectionTravelThemesUseCase struct {
	queryService selectionqueryservice.SelectionQueryService
}

func NewFetchSelectionTravelThemesUseCase(queryService selectionqueryservice.SelectionQueryService) *FetchSelectionTravelThemesUseCase {
	return &FetchSelectionTravelThemesUseCase{
		queryService: queryService,
	}
}

func (uc *FetchSelectionTravelThemesUseCase) Execute(ctx context.Context) (*api.SelectionTravelThemeResponse, error) {
	items, err := uc.queryService.FetchTravelThemes(ctx)
	if err != nil {
		return nil, err
	}

	data := make([]api.SelectionTravelThemeResponseDataInner, len(items))
	for i, item := range items {
		data[i] = api.SelectionTravelThemeResponseDataInner{
			Id:   item.ID,
			Name: item.Name,
		}
	}

	return &api.SelectionTravelThemeResponse{
		Data: data,
	}, nil
}
