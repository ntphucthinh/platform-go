package usecase

import (
	"context"

	api "platform-go/generated/api"
	homestayqueryservice "platform-go/internal/usecases/homestays/queryservice"
)

type FetchHomestaysUseCase struct {
	queryService homestayqueryservice.HomestayQueryService
}

func NewFetchHomestaysUseCase(queryService homestayqueryservice.HomestayQueryService) *FetchHomestaysUseCase {
	return &FetchHomestaysUseCase{queryService: queryService}
}

func (uc *FetchHomestaysUseCase) Execute(ctx context.Context, page, perPage int) (*api.FetchHomestayResponse, error) {
	paginated, err := uc.queryService.FetchHomestays(ctx, page, perPage)
	if err != nil {
		return nil, err
	}

	items := make([]api.FetchHomestayResponseDataInner, len(paginated.Data))
	for i, h := range paginated.Data {
		rooms := make([]api.FetchHomestayResponseDataInnerHomestayRoomsInner, len(h.Rooms))
		for j, r := range h.Rooms {
			images := make([]api.FetchHomestayResponseDataInnerHomestayRoomsInnerHomestayImagesInner, len(r.Images))
			for k, img := range r.Images {
				images[k] = api.FetchHomestayResponseDataInnerHomestayRoomsInnerHomestayImagesInner{
					Id:        img.ID,
					RoomId:    img.RoomID,
					ImagePath: img.ImagePath,
					SortOrder: img.SortOrder,
					IsCover:   img.IsCover,
				}
			}

			rooms[j] = api.FetchHomestayResponseDataInnerHomestayRoomsInner{
				Id:             r.ID,
				HomestayId:     r.HomestayID,
				Name:           r.Name,
				Description:    r.Description,
				PriceAmount:    r.PriceAmount,
				PriceUnit:      r.PriceUnit,
				SortOrder:      r.SortOrder,
				HomestayImages: images,
			}
		}

		items[i] = api.FetchHomestayResponseDataInner{
			Id:            h.ID,
			Name:          h.Name,
			Address:       h.Address,
			Description:   h.Description,
			GoogleMapsUrl: h.GoogleMapsURL,
			HomestayRooms: rooms,
		}
	}

	return &api.FetchHomestayResponse{
		Data:          items,
		PaginatorInfo: paginated.PaginatorInfo,
	}, nil
}
