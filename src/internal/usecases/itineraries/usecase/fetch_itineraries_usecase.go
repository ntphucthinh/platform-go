package usecase

import (
	"context"

	api "platform-go/generated/api"
	itineraryqueryservice "platform-go/internal/usecases/itineraries/queryservice"
)

type FetchItinerariesUseCase struct {
	queryService itineraryqueryservice.ItineraryQueryService
}

func NewFetchItinerariesUseCase(queryService itineraryqueryservice.ItineraryQueryService) *FetchItinerariesUseCase {
	return &FetchItinerariesUseCase{
		queryService: queryService,
	}
}

func (uc *FetchItinerariesUseCase) Execute(ctx context.Context, page, perPage int, locationID *int64) (*api.FetchItinerariesResponse, error) {
	paginated, err := uc.queryService.FetchItineraries(ctx, page, perPage, locationID)
	if err != nil {
		return nil, err
	}

	items := make([]api.FetchItinerariesResponseDataInner, len(paginated.Data))
	for i, it := range paginated.Data {
		items[i] = api.FetchItinerariesResponseDataInner{
			Id:             it.ID,
			LocationId:     it.LocationID,
			LocationName:   it.LocationName,
			ItineraryName:  it.ItineraryName,
			DurationDays:   it.DurationDays,
			DurationNights: it.DurationNights,
			Duration:       it.Duration,
			ActivityCount:  it.ActivityCount,
			Description:    it.Description,
			CoverImagePath: it.CoverImagePath,
		}
	}

	return &api.FetchItinerariesResponse{
		Data:          items,
		PaginatorInfo: paginated.PaginatorInfo,
	}, nil
}
