package usecase

import (
	"context"

	"gorm.io/gorm"

	api "platform-go/generated/api"
	"platform-go/internal/domain/itinerary"
	itineraryactivity "platform-go/internal/domain/itinerary_activity"
	itineraryday "platform-go/internal/domain/itinerary_day"
	itinerarytheme "platform-go/internal/domain/itinerary_theme"
	"platform-go/internal/infrastructure/repository"
)

const CreateItineraryUseCaseSuccessMessage = "Create new itinerary successfully"

type CreateItineraryUseCase struct {
	db                    repository.DBInterface
	itineraryRepo         itinerary.ItineraryRepository
	itineraryThemeRepo    itinerarytheme.ItineraryThemeRepository
	itineraryDayRepo      itineraryday.ItineraryDayRepository
	itineraryActivityRepo itineraryactivity.ItineraryActivityRepository
}

func NewCreateItineraryUseCase(
	db repository.DBInterface,
	itineraryRepo itinerary.ItineraryRepository,
	itineraryThemeRepo itinerarytheme.ItineraryThemeRepository,
	itineraryDayRepo itineraryday.ItineraryDayRepository,
	itineraryActivityRepo itineraryactivity.ItineraryActivityRepository,
) *CreateItineraryUseCase {
	return &CreateItineraryUseCase{
		db:                    db,
		itineraryRepo:         itineraryRepo,
		itineraryThemeRepo:    itineraryThemeRepo,
		itineraryDayRepo:      itineraryDayRepo,
		itineraryActivityRepo: itineraryActivityRepo,
	}
}

func (uc *CreateItineraryUseCase) Execute(
	ctx context.Context,
	request api.CreateItineraryRequest,
) (*api.MessageOnlyResponse, error) {
	err := uc.db.Transaction(ctx, func(tx *gorm.DB) error {
		txCtx := repository.WithTx(ctx, tx)

		var travelTips []string
		if request.TravelTips != nil {
			travelTips = *request.TravelTips
		}

		it, err := itinerary.NewItinerary(
			request.ItineraryName,
			request.LocationId,
			request.SuitableFor,
			request.Description,
			request.CoverImagePath,
			travelTips,
		)
		if err != nil {
			return err
		}

		if err := uc.itineraryRepo.CreateItinerary(txCtx, &it); err != nil {
			return err
		}

		if request.ThemeIds != nil && len(*request.ThemeIds) > 0 {
			themes, err := uc.buildAndValidateThemes(txCtx, it.ID, *request.ThemeIds)
			if err != nil {
				return err
			}
			if err := uc.itineraryThemeRepo.CreateItineraryThemes(txCtx, themes); err != nil {
				return err
			}
		}

		if request.Days == nil || len(*request.Days) == 0 {
			return nil
		}

		rawDays := *request.Days

		days, err := uc.buildAndValidateDays(it.ID, rawDays)
		if err != nil {
			return err
		}

		if err := uc.itineraryDayRepo.CreateItineraryDays(txCtx, days); err != nil {
			return err
		}

		activities, err := uc.buildAndValidateActivities(days, rawDays)
		if err != nil {
			return err
		}

		if len(activities) > 0 {
			if err := uc.itineraryActivityRepo.CreateItineraryActivities(txCtx, activities); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &api.MessageOnlyResponse{
		Message: CreateItineraryUseCaseSuccessMessage,
	}, nil
}

func (uc *CreateItineraryUseCase) buildAndValidateThemes(
	ctx context.Context,
	itineraryID int64,
	themeIDs []int64,
) ([]itinerarytheme.ItineraryTheme, error) {
	themes := make([]itinerarytheme.ItineraryTheme, len(themeIDs))
	seen := make(map[int64]struct{}, len(themeIDs))

	for i, themeID := range themeIDs {
		if themeID <= 0 {
			return nil, itinerarytheme.ErrInvalidThemeID
		}
		if _, exists := seen[themeID]; exists {
			return nil, itinerarytheme.ErrDuplicateThemeID
		}
		seen[themeID] = struct{}{}
		themes[i] = itinerarytheme.NewItineraryTheme(itineraryID, themeID)
	}

	validCount, err := uc.itineraryThemeRepo.CountValidThemesByIDs(ctx, themeIDs)
	if err != nil {
		return nil, err
	}
	if validCount != int64(len(themeIDs)) {
		return nil, itinerarytheme.ErrInvalidThemeID
	}

	return themes, nil
}

func (uc *CreateItineraryUseCase) buildAndValidateDays(
	itineraryID int64,
	rawDays []api.CreateItineraryDayInput,
) ([]itineraryday.ItineraryDay, error) {
	days := make([]itineraryday.ItineraryDay, len(rawDays))
	seenDayNumbers := make(map[int]struct{}, len(rawDays))

	for i, d := range rawDays {
		var dayTips []string
		if d.DayTips != nil {
			dayTips = *d.DayTips
		}

		day, err := itineraryday.NewItineraryDay(
			itineraryID,
			d.DayNumber,
			d.Title,
			d.Subtitle,
			d.Summary,
			d.ImagePath,
			dayTips,
		)
		if err != nil {
			return nil, err
		}

		if _, exists := seenDayNumbers[day.DayNumber]; exists {
			return nil, itineraryday.ErrDuplicateDayNumber
		}
		seenDayNumbers[day.DayNumber] = struct{}{}

		days[i] = day
	}

	return days, nil
}

func (uc *CreateItineraryUseCase) buildAndValidateActivities(
	days []itineraryday.ItineraryDay,
	rawDays []api.CreateItineraryDayInput,
) ([]itineraryactivity.ItineraryActivity, error) {
	var totalActivities int
	for _, d := range rawDays {
		if d.Activities != nil {
			totalActivities += len(*d.Activities)
		}
	}
	if totalActivities == 0 {
		return nil, nil
	}

	dayIDByNumber := make(map[int]int64, len(days))
	for _, day := range days {
		dayIDByNumber[day.DayNumber] = day.ID
	}

	allActivities := make([]itineraryactivity.ItineraryActivity, 0, totalActivities)

	for _, d := range rawDays {
		if d.Activities == nil || len(*d.Activities) == 0 {
			continue
		}

		dayID, exists := dayIDByNumber[d.DayNumber]
		if !exists || dayID <= 0 {
			return nil, itineraryday.ErrDayNotFound
		}

		for _, a := range *d.Activities {
			var category *string
			if a.ActivityCategory != nil {
				catStr := string(*a.ActivityCategory)
				category = &catStr
			}

			activity, err := itineraryactivity.NewItineraryActivity(
				dayID,
				a.StartTime,
				a.Title,
				category,
				a.DurationMinutes,
				a.Location,
				a.GoogleMapsUrl,
				a.Description,
				a.Notes,
				a.ImagePath,
			)
			if err != nil {
				return nil, err
			}

			allActivities = append(allActivities, activity)
		}
	}

	return allActivities, nil
}
