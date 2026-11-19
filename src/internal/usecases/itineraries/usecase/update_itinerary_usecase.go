package usecase

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	api "platform-go/generated/api"
	"platform-go/internal/domain/itinerary"
	itineraryactivity "platform-go/internal/domain/itinerary_activity"
	itineraryday "platform-go/internal/domain/itinerary_day"
	itinerarytheme "platform-go/internal/domain/itinerary_theme"
	"platform-go/internal/infrastructure/repository"
)

const UpdateItineraryUseCaseSuccessMessage = "Update itinerary successfully"

var (
	ErrUpdateItineraryInvalidID = errors.New(
		"itinerary ID must be greater than zero",
	)
	ErrUpdateItineraryInvalidDayID = errors.New(
		"day ID does not belong to this itinerary",
	)
	ErrUpdateItineraryDuplicateDayID = errors.New(
		"duplicate day ID in request",
	)
	ErrUpdateItineraryInvalidActivityID = errors.New(
		"activity ID does not belong to this day",
	)
	ErrUpdateItineraryDuplicateActivityID = errors.New(
		"duplicate activity ID in request",
	)
)

type UpdateItineraryUseCase struct {
	db                    repository.DBInterface
	itineraryRepo         itinerary.ItineraryRepository
	itineraryThemeRepo    itinerarytheme.ItineraryThemeRepository
	itineraryDayRepo      itineraryday.ItineraryDayRepository
	itineraryActivityRepo itineraryactivity.ItineraryActivityRepository
}

func NewUpdateItineraryUseCase(
	db repository.DBInterface,
	itineraryRepo itinerary.ItineraryRepository,
	itineraryThemeRepo itinerarytheme.ItineraryThemeRepository,
	itineraryDayRepo itineraryday.ItineraryDayRepository,
	itineraryActivityRepo itineraryactivity.ItineraryActivityRepository,
) *UpdateItineraryUseCase {
	return &UpdateItineraryUseCase{
		db:                    db,
		itineraryRepo:         itineraryRepo,
		itineraryThemeRepo:    itineraryThemeRepo,
		itineraryDayRepo:      itineraryDayRepo,
		itineraryActivityRepo: itineraryActivityRepo,
	}
}

type pendingNewActivity struct {
	dayNumber      int
	itineraryDayID int64
	activity       itineraryactivity.ItineraryActivity
}

type itinerarySyncPlan struct {
	dayIDsToDelete            []int64
	daysToUpdate              []itineraryday.ItineraryDay
	daysToCreate              []itineraryday.ItineraryDay
	activityIDsToDelete       []int64
	activitiesToUpdate        []itineraryactivity.ItineraryActivity
	pendingActivitiesToCreate []pendingNewActivity
}

func (uc *UpdateItineraryUseCase) Execute(
	ctx context.Context,
	itineraryID int64,
	request api.UpdateItineraryRequest,
) (*api.MessageOnlyResponse, error) {
	if itineraryID <= 0 {
		return nil, ErrUpdateItineraryInvalidID
	}

	err := uc.db.Transaction(ctx, func(tx *gorm.DB) error {
		txCtx := repository.WithTx(ctx, tx)

		existingItinerary, err := uc.itineraryRepo.FindItineraryByID(txCtx, itineraryID)
		if err != nil {
			return err
		}
		if existingItinerary == nil {
			return ErrItineraryNotFound
		}

		var travelTips []string
		if request.TravelTips != nil {
			travelTips = *request.TravelTips
		}

		if err := existingItinerary.Update(
			request.ItineraryName,
			request.LocationId,
			request.SuitableFor,
			request.Description,
			request.CoverImagePath,
			travelTips,
		); err != nil {
			return err
		}

		var themesToCreate []itinerarytheme.ItineraryTheme
		if request.ThemeIds != nil && len(*request.ThemeIds) > 0 {
			themesToCreate, err = uc.buildAndValidateThemes(txCtx, existingItinerary.ID, *request.ThemeIds)
			if err != nil {
				return err
			}
		}

		var plan *itinerarySyncPlan
		if request.Days != nil {
			plan, err = uc.buildAndValidateSyncPlan(txCtx, existingItinerary.ID, *request.Days)
			if err != nil {
				return err
			}
		}

		if err := uc.itineraryRepo.UpdateItinerary(txCtx, existingItinerary); err != nil {
			return err
		}

		if request.ThemeIds != nil {
			if err := uc.itineraryThemeRepo.DeleteThemesByItineraryID(txCtx, existingItinerary.ID); err != nil {
				return err
			}
			if len(themesToCreate) > 0 {
				if err := uc.itineraryThemeRepo.CreateItineraryThemes(txCtx, themesToCreate); err != nil {
					return err
				}
			}
		}

		if plan != nil {
			if err := uc.applySyncPlan(txCtx, plan); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &api.MessageOnlyResponse{
		Message: UpdateItineraryUseCaseSuccessMessage,
	}, nil
}

func (uc *UpdateItineraryUseCase) buildAndValidateThemes(
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

func (uc *UpdateItineraryUseCase) buildAndValidateSyncPlan(
	txCtx context.Context,
	itineraryID int64,
	rawDays []api.UpdateItineraryDayInput,
) (*itinerarySyncPlan, error) {
	existingDays, err := uc.itineraryDayRepo.FindDaysByItineraryID(txCtx, itineraryID)
	if err != nil {
		return nil, err
	}

	existingDayMap := make(map[int64]itineraryday.ItineraryDay, len(existingDays))
	existingDayIDs := make([]int64, len(existingDays))
	for i, day := range existingDays {
		existingDayMap[day.ID] = day
		existingDayIDs[i] = day.ID
	}

	if len(rawDays) == 0 {
		return &itinerarySyncPlan{
			dayIDsToDelete: existingDayIDs,
		}, nil
	}

	var existingActivities []itineraryactivity.ItineraryActivity
	if len(existingDayIDs) > 0 {
		existingActivities, err = uc.itineraryActivityRepo.FindActivitiesByDayIDs(txCtx, existingDayIDs)
		if err != nil {
			return nil, err
		}
	}

	existingActivityMapByDayID := make(map[int64]map[int64]itineraryactivity.ItineraryActivity, len(existingDays))
	existingActivityListByDayID := make(map[int64][]itineraryactivity.ItineraryActivity, len(existingDays))
	for _, act := range existingActivities {
		if existingActivityMapByDayID[act.ItineraryDayID] == nil {
			existingActivityMapByDayID[act.ItineraryDayID] = make(map[int64]itineraryactivity.ItineraryActivity)
		}
		existingActivityMapByDayID[act.ItineraryDayID][act.ID] = act
		existingActivityListByDayID[act.ItineraryDayID] = append(existingActivityListByDayID[act.ItineraryDayID], act)
	}

	plan := &itinerarySyncPlan{
		daysToUpdate: make([]itineraryday.ItineraryDay, 0, len(rawDays)),
		daysToCreate: make([]itineraryday.ItineraryDay, 0, len(rawDays)),
	}

	requestedDayIDs := make(map[int64]struct{}, len(rawDays))
	requestedActivityIDs := make(map[int64]struct{})
	seenDayNumbers := make(map[int]struct{}, len(rawDays))

	for _, dayReq := range rawDays {
		var isExistingDay bool
		var existingDayID int64

		if dayReq.Id != nil {
			dayID := *dayReq.Id
			if dayID <= 0 {
				return nil, fmt.Errorf("%w: %d", ErrUpdateItineraryInvalidDayID, dayID)
			}
			if _, exists := requestedDayIDs[dayID]; exists {
				return nil, fmt.Errorf("%w: %d", ErrUpdateItineraryDuplicateDayID, dayID)
			}
			if _, exists := existingDayMap[dayID]; !exists {
				return nil, fmt.Errorf("%w: %d", ErrUpdateItineraryInvalidDayID, dayID)
			}
			requestedDayIDs[dayID] = struct{}{}
			isExistingDay = true
			existingDayID = dayID
		}

		var dayTips []string
		if dayReq.DayTips != nil {
			dayTips = *dayReq.DayTips
		}

		dayEntity, err := itineraryday.NewItineraryDay(
			itineraryID,
			dayReq.DayNumber,
			dayReq.Title,
			dayReq.Subtitle,
			dayReq.Summary,
			dayReq.ImagePath,
			dayTips,
		)
		if err != nil {
			return nil, err
		}

		if _, exists := seenDayNumbers[dayEntity.DayNumber]; exists {
			return nil, itineraryday.ErrDuplicateDayNumber
		}
		seenDayNumbers[dayEntity.DayNumber] = struct{}{}

		if isExistingDay {
			dayEntity.ID = existingDayID
			plan.daysToUpdate = append(plan.daysToUpdate, dayEntity)
		} else {
			plan.daysToCreate = append(plan.daysToCreate, dayEntity)
		}

		if dayReq.Activities == nil {
			continue
		}

		if err := validateAndClassifyDayActivities(
			*dayReq.Activities,
			isExistingDay,
			existingDayID,
			dayEntity.DayNumber,
			existingActivityMapByDayID[existingDayID],
			existingActivityListByDayID[existingDayID],
			requestedActivityIDs,
			plan,
		); err != nil {
			return nil, err
		}
	}

	for _, existingDay := range existingDays {
		if _, exists := requestedDayIDs[existingDay.ID]; !exists {
			plan.dayIDsToDelete = append(plan.dayIDsToDelete, existingDay.ID)
		}
	}

	return plan, nil
}

func validateAndClassifyDayActivities(
	rawActivities []api.UpdateItineraryActivityInput,
	isExistingDay bool,
	existingDayID int64,
	dayNumber int,
	dayExistingActivityMap map[int64]itineraryactivity.ItineraryActivity,
	dayExistingActivityList []itineraryactivity.ItineraryActivity,
	requestedActivityIDs map[int64]struct{},
	plan *itinerarySyncPlan,
) error {
	dayRequestedActivityIDs := make(map[int64]struct{}, len(rawActivities))

	for _, actReq := range rawActivities {
		var isExistingActivity bool
		var existingActivityID int64

		if actReq.Id != nil {
			actID := *actReq.Id
			if actID <= 0 {
				return fmt.Errorf("%w: %d", ErrUpdateItineraryInvalidActivityID, actID)
			}
			if _, exists := requestedActivityIDs[actID]; exists {
				return fmt.Errorf("%w: %d", ErrUpdateItineraryDuplicateActivityID, actID)
			}
			if !isExistingDay {
				return fmt.Errorf("%w: %d", ErrUpdateItineraryInvalidActivityID, actID)
			}
			if _, exists := dayExistingActivityMap[actID]; !exists {
				return fmt.Errorf("%w: %d", ErrUpdateItineraryInvalidActivityID, actID)
			}
			requestedActivityIDs[actID] = struct{}{}
			dayRequestedActivityIDs[actID] = struct{}{}
			isExistingActivity = true
			existingActivityID = actID
		}

		var category *string
		if actReq.ActivityCategory != nil {
			catStr := string(*actReq.ActivityCategory)
			category = &catStr
		}

		actEntity, err := itineraryactivity.NewItineraryActivity(
			existingDayID,
			actReq.StartTime,
			actReq.Title,
			category,
			actReq.DurationMinutes,
			actReq.Location,
			actReq.GoogleMapsUrl,
			actReq.Description,
			actReq.Notes,
			actReq.ImagePath,
		)
		if err != nil {
			return err
		}

		if isExistingActivity {
			actEntity.ID = existingActivityID
			plan.activitiesToUpdate = append(plan.activitiesToUpdate, actEntity)
		} else {
			plan.pendingActivitiesToCreate = append(plan.pendingActivitiesToCreate, pendingNewActivity{
				dayNumber:      dayNumber,
				itineraryDayID: existingDayID,
				activity:       actEntity,
			})
		}
	}

	if isExistingDay {
		for _, existingAct := range dayExistingActivityList {
			if _, exists := dayRequestedActivityIDs[existingAct.ID]; !exists {
				plan.activityIDsToDelete = append(plan.activityIDsToDelete, existingAct.ID)
			}
		}
	}

	return nil
}

func (uc *UpdateItineraryUseCase) applySyncPlan(
	txCtx context.Context,
	plan *itinerarySyncPlan,
) error {
	if len(plan.dayIDsToDelete) > 0 {
		if err := uc.itineraryActivityRepo.DeleteActivitiesByDayIDs(txCtx, plan.dayIDsToDelete); err != nil {
			return err
		}
		if err := uc.itineraryDayRepo.DeleteDaysByIDs(txCtx, plan.dayIDsToDelete); err != nil {
			return err
		}
	}

	if len(plan.daysToUpdate) > 0 {
		if err := uc.itineraryDayRepo.UpdateItineraryDays(txCtx, plan.daysToUpdate); err != nil {
			return err
		}
	}

	if len(plan.daysToCreate) > 0 {
		if err := uc.itineraryDayRepo.CreateItineraryDays(txCtx, plan.daysToCreate); err != nil {
			return err
		}
	}

	if len(plan.activityIDsToDelete) > 0 {
		if err := uc.itineraryActivityRepo.DeleteActivitiesByIDs(txCtx, plan.activityIDsToDelete); err != nil {
			return err
		}
	}

	if len(plan.activitiesToUpdate) > 0 {
		if err := uc.itineraryActivityRepo.UpdateItineraryActivities(txCtx, plan.activitiesToUpdate); err != nil {
			return err
		}
	}

	if len(plan.pendingActivitiesToCreate) > 0 {
		createdDayIDByNumber := make(map[int]int64, len(plan.daysToCreate))
		for _, d := range plan.daysToCreate {
			createdDayIDByNumber[d.DayNumber] = d.ID
		}

		activitiesToCreate := make([]itineraryactivity.ItineraryActivity, len(plan.pendingActivitiesToCreate))
		for i, pending := range plan.pendingActivitiesToCreate {
			dayID := pending.itineraryDayID
			if dayID <= 0 {
				var exists bool
				dayID, exists = createdDayIDByNumber[pending.dayNumber]
				if !exists || dayID <= 0 {
					return itineraryday.ErrDayNotFound
				}
			}

			act := pending.activity
			act.ItineraryDayID = dayID
			activitiesToCreate[i] = act
		}

		if err := uc.itineraryActivityRepo.CreateItineraryActivities(txCtx, activitiesToCreate); err != nil {
			return err
		}
	}

	return nil
}
