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

func (uc *FetchHomestaysUseCase) Execute(ctx context.Context) ([]api.FetchHomestayResponse, error) {
	homestays, err := uc.queryService.FetchHomestays(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]api.FetchHomestayResponse, len(homestays))
	for i, h := range homestays {
		rooms := make([]api.FetchHomestayResponseHomestayRoomsInner, len(h.Rooms))
		for j, r := range h.Rooms {
			images := make([]api.FetchHomestayResponseHomestayRoomsInnerHomestayImagesInner, len(r.Images))
			for k, img := range r.Images {
				images[k] = api.FetchHomestayResponseHomestayRoomsInnerHomestayImagesInner{
					Id:        img.ID,
					RoomId:    img.RoomID,
					ImagePath: img.ImagePath,
					SortOrder: img.SortOrder,
					IsCover:   img.IsCover,
				}
			}

			rooms[j] = api.FetchHomestayResponseHomestayRoomsInner{
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

		items[i] = api.FetchHomestayResponse{
			Id:            h.ID,
			Name:          h.Name,
			Address:       h.Address,
			Description:   h.Description,
			GoogleMapsUrl: h.GoogleMapsURL,
			HomestayRooms: rooms,
		}
	}

	return items, nil
}
