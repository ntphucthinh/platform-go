package usecase

import (
	"context"
	"errors"

	openapi_types "github.com/oapi-codegen/runtime/types"

	api "platform-go/generated/api"
	travelmemoryqueryservice "platform-go/internal/usecases/travel_memories/queryservice"
)

var ErrTravelMemoryNotFound = errors.New("travel memory not found")

type FetchTravelMemoryByIDUseCase struct {
	queryService travelmemoryqueryservice.TravelMemoryQueryService
}

func NewFetchTravelMemoryByIDUseCase(queryService travelmemoryqueryservice.TravelMemoryQueryService) *FetchTravelMemoryByIDUseCase {
	return &FetchTravelMemoryByIDUseCase{
		queryService: queryService,
	}
}

func (uc *FetchTravelMemoryByIDUseCase) Execute(ctx context.Context, id int64) (*api.FetchTravelMemoryByIDResponse, error) {
	m, err := uc.queryService.FetchTravelMemoryByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, ErrTravelMemoryNotFound
	}

	return &api.FetchTravelMemoryByIDResponse{
		Data: api.FetchTravelMemoryByIDResponseData{
			Id:             m.ID,
			LocationId:     m.LocationID,
			Title:          m.Title,
			StartDate:      openapi_types.Date{Time: m.StartDate},
			EndDate:        openapi_types.Date{Time: m.EndDate},
			Description:    m.Description,
			CoverImagePath: m.CoverImagePath,
			Location: api.FetchTravelMemoryByIDResponseDataLocation{
				Id:           m.LocationID,
				LocationName: m.LocationName,
				CountryCode:  m.CountryCode,
				CountryName:  m.CountryName,
			},
			Days: mapTravelMemoryDays(m.Days),
		},
	}, nil
}
