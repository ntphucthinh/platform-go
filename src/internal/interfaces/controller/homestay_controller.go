package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	api "platform-go/generated/api"
	"platform-go/internal/common"
	"platform-go/internal/infrastructure/queryservice"
)

type HomestayController struct {
	queryService *queryservice.HomestayQueryService
}

func NewHomestayController(queryService *queryservice.HomestayQueryService) *HomestayController {
	return &HomestayController{
		queryService: queryService,
	}
}

// Get homestay list
// (GET /homestays)
func (hc *HomestayController) FetchHomestays(c *gin.Context) {
	homestays, err := hc.queryService.FetchHomestays(c.Request.Context())
	if err != nil {
		common.IsErrorWithMessage(c, err, "Failed to fetch homestays")
		return
	}

	items := make([]api.FetchHomestayResponse, len(homestays))
	for i, h := range homestays {
		rooms := make([]api.FetchHomestayResponseHomestayRoomsInner, len(h.Rooms))
		for j, r := range h.Rooms {
			images := make([]api.FetchHomestayResponseHomestayRoomsInnerHomestayImagesInner, len(r.Images))
			for k, img := range r.Images {
				images[k] = api.FetchHomestayResponseHomestayRoomsInnerHomestayImagesInner{
					Id:         img.ID,
					HomestayId: img.HomestayID,
					RoomId:     img.RoomID,
					ImagePath:  img.ImagePath,
					SortOrder:  img.SortOrder,
					IsCover:    img.IsCover,
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

	c.JSON(http.StatusOK, api.FetchHomestays200Response{
		Data: items,
	})
}
