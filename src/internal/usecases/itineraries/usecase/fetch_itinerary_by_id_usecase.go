package usecase

import (
	"context"
	"errors"

	api "platform-go/generated/api"
	"platform-go/internal/infrastructure/qsdto"
	itineraryqueryservice "platform-go/internal/usecases/itineraries/queryservice"
)

var ErrItineraryNotFound = errors.New("itinerary not found")

type FetchItineraryByIDUseCase struct {
	queryService itineraryqueryservice.ItineraryQueryService
}

func NewFetchItineraryByIDUseCase(queryService itineraryqueryservice.ItineraryQueryService) *FetchItineraryByIDUseCase {
	return &FetchItineraryByIDUseCase{
		queryService: queryService,
	}
}

func (uc *FetchItineraryByIDUseCase) Execute(ctx context.Context, id int64) (*api.FetchItineraryByIDResponse, error) {
	it, err := uc.queryService.FetchItineraryByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if it == nil {
		return nil, ErrItineraryNotFound
	}

	themes := make([]api.FetchItineraryByIDResponseDataThemesInner, len(it.Themes))
	for i, th := range it.Themes {
		themes[i] = api.FetchItineraryByIDResponseDataThemesInner{
			Id:   th.ID,
			Name: th.Name,
		}
	}

	featuredImages := make([]api.FetchItineraryByIDResponseDataFeaturedImagesInner, len(it.FeaturedImages))
	for i, img := range it.FeaturedImages {
		featuredImages[i] = api.FetchItineraryByIDResponseDataFeaturedImagesInner{
			Id:        img.ID,
			ImagePath: img.ImagePath,
		}
	}

	travelTips := it.TravelTips
	if travelTips == nil {
		travelTips = []string{}
	}

	return &api.FetchItineraryByIDResponse{
		Data: api.FetchItineraryByIDResponseData{
			Id:             it.ID,
			ItineraryName:  it.ItineraryName,
			LocationId:     it.LocationID,
			LocationName:   it.LocationName,
			DurationDays:   it.DurationDays,
			DurationNights: it.DurationNights,
			Duration:       it.Duration,
			SuitableFor:    it.SuitableFor,
			ActivityCount:  it.ActivityCount,
			Themes:         themes,
			Description:    it.Description,
			CoverImagePath: it.CoverImagePath,
			TravelTips:     travelTips,
			FeaturedImages: featuredImages,
			Days:           mapItineraryDays(it.Days),
		},
	}, nil
}

func mapItineraryDays(dtos []qsdto.ItineraryDayDto) []api.FetchItineraryByIDResponseDataDaysInner {
	days := make([]api.FetchItineraryByIDResponseDataDaysInner, len(dtos))
	for i, d := range dtos {
		dayTips := d.DayTips
		if dayTips == nil {
			dayTips = []string{}
		}

		activities := make([]api.FetchItineraryByIDResponseDataDaysInnerActivitiesInner, len(d.Activities))
		for j, a := range d.Activities {
			var category *api.ItineraryActivityCategory
			if a.ActivityCategory != nil {
				cat := api.ItineraryActivityCategory(*a.ActivityCategory)
				category = &cat
			}

			activities[j] = api.FetchItineraryByIDResponseDataDaysInnerActivitiesInner{
				Id:               a.ID,
				StartTime:        a.StartTime,
				Title:            a.Title,
				ActivityCategory: category,
				DurationMinutes:  a.DurationMinutes,
				Location:         a.Location,
				GoogleMapsUrl:    a.GoogleMapsURL,
				Description:      a.Description,
				Notes:            a.Notes,
				ImagePath:        a.ImagePath,
			}
		}

		days[i] = api.FetchItineraryByIDResponseDataDaysInner{
			Id:         d.ID,
			DayNumber:  d.DayNumber,
			Title:      d.Title,
			Subtitle:   d.Subtitle,
			Summary:    d.Summary,
			ImagePath:  d.ImagePath,
			DayTips:    dayTips,
			Activities: activities,
		}
	}
	return days
}
