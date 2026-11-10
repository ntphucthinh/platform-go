package usecase

import (
	"context"

	openapi_types "github.com/oapi-codegen/runtime/types"

	api "platform-go/generated/api"
	"platform-go/internal/infrastructure/qsdto"
	travelmemoryqueryservice "platform-go/internal/usecases/travel_memories/queryservice"
)

const defaultTravelMemoriesPerPage = 10

type FetchTravelMemoriesUseCase struct {
	queryService travelmemoryqueryservice.TravelMemoryQueryService
}

func NewFetchTravelMemoriesUseCase(queryService travelmemoryqueryservice.TravelMemoryQueryService) *FetchTravelMemoriesUseCase {
	return &FetchTravelMemoriesUseCase{
		queryService: queryService,
	}
}

func (uc *FetchTravelMemoriesUseCase) Execute(ctx context.Context, page int, countryCode *string, locationID *int64) (*api.FetchTravelMemoriesResponse, error) {
	paginated, err := uc.queryService.FetchTravelMemories(ctx, page, defaultTravelMemoriesPerPage, countryCode, locationID)
	if err != nil {
		return nil, err
	}

	items := make([]api.FetchTravelMemoriesResponseDataInner, len(paginated.Data))
	for i, m := range paginated.Data {
		items[i] = api.FetchTravelMemoriesResponseDataInner{
			Id:             m.ID,
			LocationId:     m.LocationID,
			LocationName:   m.LocationName,
			CountryCode:    m.CountryCode,
			Title:          m.Title,
			StartDate:      openapi_types.Date{Time: m.StartDate},
			EndDate:        openapi_types.Date{Time: m.EndDate},
			Description:    m.Description,
			CoverImagePath: m.CoverImagePath,
			Days:           mapTravelMemoryDays(m.Days),
		}
	}

	return &api.FetchTravelMemoriesResponse{
		Data:          items,
		PaginatorInfo: paginated.PaginatorInfo,
	}, nil
}

func mapTravelMemoryDays(days []qsdto.TravelMemoryDayDto) []api.FetchTravelMemoriesResponseDataInnerDaysInner {
	result := make([]api.FetchTravelMemoriesResponseDataInnerDaysInner, len(days))
	for i, d := range days {
		photos := make([]api.FetchTravelMemoriesResponseDataInnerDaysInnerPhotosInner, len(d.Photos))
		for j, p := range d.Photos {
			photos[j] = api.FetchTravelMemoriesResponseDataInnerDaysInnerPhotosInner{
				Id:          p.ID,
				MemoryDayId: p.MemoryDayID,
				ImagePath:   p.ImagePath,
				SortOrder:   p.SortOrder,
			}
		}

		result[i] = api.FetchTravelMemoriesResponseDataInnerDaysInner{
			Id:          d.ID,
			MemoryId:    d.MemoryID,
			Date:        openapi_types.Date{Time: d.Date},
			Description: d.Description,
			Photos:      photos,
		}
	}

	return result
}
