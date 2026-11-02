package usecase

import (
	"context"
	"errors"

	api "platform-go/generated/api"
	homestayqueryservice "platform-go/internal/usecases/homestays/queryservice"
)

var ErrHomestayNotFound = errors.New("homestay not found")

type FetchHomestayByIDUseCase struct {
	queryService homestayqueryservice.HomestayQueryService
}

func NewFetchHomestayByIDUseCase(queryService homestayqueryservice.HomestayQueryService) *FetchHomestayByIDUseCase {
	return &FetchHomestayByIDUseCase{queryService: queryService}
}

func (uc *FetchHomestayByIDUseCase) Execute(ctx context.Context, id int64) (*api.FetchHomestayByIDResponse, error) {
	h, err := uc.queryService.FetchHomestayByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if h == nil {
		return nil, ErrHomestayNotFound
	}

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
			PriceUnit:      (*api.PriceUnit)(r.PriceUnit),
			SortOrder:      r.SortOrder,
			HomestayImages: images,
		}
	}

	return &api.FetchHomestayByIDResponse{
		Data: api.FetchHomestayResponseDataInner{
			Id:            h.ID,
			Name:          h.Name,
			Address:       h.Address,
			Description:   h.Description,
			GoogleMapsUrl: h.GoogleMapsURL,
			HomestayRooms: rooms,
		},
	}, nil
}
