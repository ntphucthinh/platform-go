package itinerary_controller_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	api "platform-go/generated/api"
	"platform-go/internal/controller_test/shared"
	"platform-go/internal/domain/itinerary"
	itineraryactivity "platform-go/internal/domain/itinerary_activity"
	itineraryday "platform-go/internal/domain/itinerary_day"
	itinerarytheme "platform-go/internal/domain/itinerary_theme"
	"platform-go/internal/infrastructure/dbmodel"
	itineraryusecase "platform-go/internal/usecases/itineraries/usecase"
	"platform-go/tests/factory"
)

func TestItineraryController_UpdateItinerary(t *testing.T) {
	t.Run("Normal_UpdateItinerarySuccess", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		loc1 := factory.NewTravelLocationFactory().
			WithLocationName("Đà Lạt").
			WithCountryCode("VN").
			WithCountryName("Vietnam").
			MustInsert(shared.TestDB)

		loc2 := factory.NewTravelLocationFactory().
			WithLocationName("Nha Trang").
			WithCountryCode("VN").
			WithCountryName("Vietnam").
			MustInsert(shared.TestDB)

		theme1 := factory.NewTravelThemeFactory().
			WithName("Thiên nhiên").
			MustInsert(shared.TestDB)
		theme2 := factory.NewTravelThemeFactory().
			WithName("Ẩm thực").
			MustInsert(shared.TestDB)
		theme3 := factory.NewTravelThemeFactory().
			WithName("Biển đảo").
			MustInsert(shared.TestDB)

		initialItinerary := factory.NewItineraryFactory().
			WithLocationID(&loc1.ID).
			WithName("Khám phá Đà Lạt ban đầu").
			WithSuitableFor("Cặp đôi").
			WithDescription("Mô tả ban đầu").
			WithCoverImagePath("itineraries/cover-old.jpg").
			WithTravelTips(`["Tip cũ 1"]`).
			MustInsert(shared.TestDB)

		factory.NewItineraryThemeFactory().
			WithItineraryID(initialItinerary.ID).
			WithThemeID(theme1.ID).
			MustInsert(shared.TestDB)

		day1 := factory.NewItineraryDayFactory().
			WithItineraryID(initialItinerary.ID).
			WithDayNumber(1).
			WithTitle("Ngày 1 cũ").
			WithImagePath("itineraries/days/day1-old.jpg").
			MustInsert(shared.TestDB)

		day2 := factory.NewItineraryDayFactory().
			WithItineraryID(initialItinerary.ID).
			WithDayNumber(2).
			WithTitle("Ngày 2 sẽ bị xóa").
			MustInsert(shared.TestDB)

		act1 := factory.NewItineraryActivityFactory().
			WithItineraryDayID(day1.ID).
			WithStartTime("07:30:00").
			WithTitle("Ăn sáng cũ").
			MustInsert(shared.TestDB)

		act2 := factory.NewItineraryActivityFactory().
			WithItineraryDayID(day1.ID).
			WithStartTime("09:00:00").
			WithTitle("Hoạt động sẽ bị xóa").
			MustInsert(shared.TestDB)

		act3 := factory.NewItineraryActivityFactory().
			WithItineraryDayID(day2.ID).
			WithStartTime("10:00:00").
			WithTitle("Hoạt động của ngày 2 sẽ bị xóa").
			MustInsert(shared.TestDB)

		newSuitableFor := "Gia đình, nhóm bạn"
		newDesc := "Hành trình Nha Trang đã cập nhật"
		newCover := "itineraries/cover-updated.jpg"
		newTips := []string{"Mang kem chống nắng", "Chuẩn bị đồ bơi"}
		newThemes := []int64{theme2.ID, theme3.ID}

		day1Subtitle := "Phụ đề ngày 1 mới"
		day1Summary := "Tóm tắt ngày 1 mới"
		day1Img := "itineraries/days/day1-updated.jpg"
		day1Tips := []string{"Đi sớm"}

		actCategoryDining := api.ItineraryActivityCategoryDining
		actCategoryRest := api.ItineraryActivityCategoryRest
		act1Duration := 45
		act1Location := "Quán mới"
		act1Img := "itineraries/activities/act1-updated.jpg"
		newActImg := "itineraries/activities/act-new.jpg"

		payload := api.UpdateItineraryRequest{
			ItineraryName:  "Khám phá Nha Trang cập nhật",
			LocationId:     &loc2.ID,
			SuitableFor:    &newSuitableFor,
			Description:    &newDesc,
			CoverImagePath: &newCover,
			TravelTips:     &newTips,
			ThemeIds:       &newThemes,
			Days: &[]api.UpdateItineraryDayInput{
				{
					Id:        &day1.ID,
					DayNumber: 1,
					Title:     "Ngày 1 đã cập nhật",
					Subtitle:  &day1Subtitle,
					Summary:   &day1Summary,
					ImagePath: &day1Img,
					DayTips:   &day1Tips,
					Activities: &[]api.UpdateItineraryActivityInput{
						{
							Id:               &act1.ID,
							StartTime:        "08:00",
							Title:            "Ăn sáng đã cập nhật",
							ActivityCategory: &actCategoryDining,
							DurationMinutes:  &act1Duration,
							Location:         &act1Location,
							ImagePath:        &act1Img,
						},
						{
							StartTime:        "14:00",
							Title:            "Nghỉ trưa tại khách sạn",
							ActivityCategory: &actCategoryRest,
							ImagePath:        &newActImg,
						},
					},
				},
				{
					DayNumber: 3,
					Title:     "Ngày 3 mới tạo",
					Activities: &[]api.UpdateItineraryActivityInput{
						{
							StartTime: "09:30",
							Title:     "Tắm biển ngày 3",
						},
					},
				},
			},
		}

		req := shared.NewRequest(t, admin, http.MethodPut, fmt.Sprintf("/itineraries/%d", initialItinerary.ID), payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())

		var res api.MessageOnlyResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
		assert.Equal(t, itineraryusecase.UpdateItineraryUseCaseSuccessMessage, res.Message)

		// Verify itinerary updated
		var updatedIt dbmodel.Itinerary
		require.NoError(t, shared.TestDB.First(&updatedIt, initialItinerary.ID).Error)
		assert.Equal(t, "Khám phá Nha Trang cập nhật", updatedIt.ItineraryName)
		require.NotNil(t, updatedIt.LocationID)
		assert.Equal(t, loc2.ID, *updatedIt.LocationID)
		require.NotNil(t, updatedIt.SuitableFor)
		assert.Equal(t, newSuitableFor, *updatedIt.SuitableFor)
		require.NotNil(t, updatedIt.Description)
		assert.Equal(t, newDesc, *updatedIt.Description)
		require.NotNil(t, updatedIt.CoverImagePath)
		assert.Equal(t, newCover, *updatedIt.CoverImagePath)

		// Verify themes replaced
		var themes []dbmodel.ItineraryTheme
		require.NoError(t, shared.TestDB.Where("itinerary_id = ?", initialItinerary.ID).Order("theme_id ASC").Find(&themes).Error)
		require.Len(t, themes, 2)
		assert.Equal(t, theme2.ID, themes[0].ThemeID)
		assert.Equal(t, theme3.ID, themes[1].ThemeID)

		// Verify Day 1 updated
		var updatedDay1 dbmodel.ItineraryDay
		require.NoError(t, shared.TestDB.First(&updatedDay1, day1.ID).Error)
		assert.Equal(t, "Ngày 1 đã cập nhật", updatedDay1.Title)
		require.NotNil(t, updatedDay1.Subtitle)
		assert.Equal(t, day1Subtitle, *updatedDay1.Subtitle)
		require.NotNil(t, updatedDay1.ImagePath)
		assert.Equal(t, day1Img, *updatedDay1.ImagePath)

		// Verify Day 2 deleted
		var deletedDay2 dbmodel.ItineraryDay
		assert.True(t, errors.Is(shared.TestDB.First(&deletedDay2, day2.ID).Error, gorm.ErrRecordNotFound))

		// Verify Day 1 activities (act1 updated, new activity created, act2 deleted)
		var day1Acts []dbmodel.ItineraryActivity
		require.NoError(t, shared.TestDB.Where("itinerary_day_id = ?", day1.ID).Order("start_time ASC").Find(&day1Acts).Error)
		require.Len(t, day1Acts, 2)
		assert.Equal(t, act1.ID, day1Acts[0].ID)
		assert.Equal(t, "08:00:00", day1Acts[0].StartTime)
		assert.Equal(t, "Ăn sáng đã cập nhật", day1Acts[0].Title)
		assert.Equal(t, "14:00:00", day1Acts[1].StartTime)
		assert.Equal(t, "Nghỉ trưa tại khách sạn", day1Acts[1].Title)

		var deletedAct dbmodel.ItineraryActivity
		assert.True(t, errors.Is(shared.TestDB.First(&deletedAct, act2.ID).Error, gorm.ErrRecordNotFound))
		assert.True(t, errors.Is(shared.TestDB.First(&deletedAct, act3.ID).Error, gorm.ErrRecordNotFound))

		// Verify newly created Day 3 and its activity
		var createdDay3 dbmodel.ItineraryDay
		require.NoError(t, shared.TestDB.Where("itinerary_id = ? AND day_number = 3", initialItinerary.ID).First(&createdDay3).Error)
		assert.Equal(t, "Ngày 3 mới tạo", createdDay3.Title)

		var day3Acts []dbmodel.ItineraryActivity
		require.NoError(t, shared.TestDB.Where("itinerary_day_id = ?", createdDay3.ID).Find(&day3Acts).Error)
		require.Len(t, day3Acts, 1)
		assert.Equal(t, "Tắm biển ngày 3", day3Acts[0].Title)
	})

	t.Run("Normal_ThemesNil_EmptyThemes_DaysNil_EmptyDays_ActivitiesNil_EmptyActivities_ZeroValues", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		loc := factory.NewTravelLocationFactory().
			WithLocationName("Đà Nẵng").
			WithCountryCode("VN").
			MustInsert(shared.TestDB)

		theme := factory.NewTravelThemeFactory().
			WithName("Thiên nhiên").
			MustInsert(shared.TestDB)

		it := factory.NewItineraryFactory().
			WithLocationID(&loc.ID).
			WithName("Initial Title").
			WithSuitableFor("Gia đình").
			WithDescription("Initial Description").
			WithCoverImagePath("itineraries/initial-cover.jpg").
			WithTravelTips(`["Tip 1"]`).
			MustInsert(shared.TestDB)

		factory.NewItineraryThemeFactory().
			WithItineraryID(it.ID).
			WithThemeID(theme.ID).
			MustInsert(shared.TestDB)

		day1 := factory.NewItineraryDayFactory().
			WithItineraryID(it.ID).
			WithDayNumber(1).
			WithTitle("Day 1").
			WithSubtitle("Sub 1").
			MustInsert(shared.TestDB)

		day2 := factory.NewItineraryDayFactory().
			WithItineraryID(it.ID).
			WithDayNumber(2).
			WithTitle("Day 2").
			MustInsert(shared.TestDB)

		act1 := factory.NewItineraryActivityFactory().
			WithItineraryDayID(day1.ID).
			WithStartTime("08:00:00").
			WithTitle("Act 1").
			MustInsert(shared.TestDB)

		act2 := factory.NewItineraryActivityFactory().
			WithItineraryDayID(day2.ID).
			WithStartTime("09:00:00").
			WithTitle("Act 2").
			MustInsert(shared.TestDB)

		// 1. ThemeIds == nil and Days == nil: keeps existing themes, days, and activities while clearing optional parent fields
		payloadNil := api.UpdateItineraryRequest{
			ItineraryName:  "Updated Name Only",
			LocationId:     nil,
			SuitableFor:    nil,
			Description:    nil,
			CoverImagePath: nil,
			TravelTips:     nil,
			ThemeIds:       nil,
			Days:           nil,
		}
		req1 := shared.NewRequest(t, admin, http.MethodPut, fmt.Sprintf("/itineraries/%d", it.ID), payloadNil)
		w1 := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w1, req1)
		require.Equal(t, http.StatusOK, w1.Code, w1.Body.String())

		var itAfter1 dbmodel.Itinerary
		require.NoError(t, shared.TestDB.First(&itAfter1, it.ID).Error)
		assert.Nil(t, itAfter1.LocationID)
		assert.Nil(t, itAfter1.SuitableFor)
		assert.Nil(t, itAfter1.Description)
		assert.Nil(t, itAfter1.CoverImagePath)
		assert.Nil(t, itAfter1.TravelTips)

		var themeCount1, dayCount1 int64
		shared.TestDB.Model(&dbmodel.ItineraryTheme{}).Where("itinerary_id = ?", it.ID).Count(&themeCount1)
		shared.TestDB.Model(&dbmodel.ItineraryDay{}).Where("itinerary_id = ?", it.ID).Count(&dayCount1)
		assert.Equal(t, int64(1), themeCount1)
		assert.Equal(t, int64(2), dayCount1)

		// 2. Activities == nil on Day 1 (keeps act1), Empty Activities on Day 2 (deletes act2), Empty ThemeIds (deletes themes)
		emptyThemes := []int64{}
		emptyActs := []api.UpdateItineraryActivityInput{}
		payloadSync := api.UpdateItineraryRequest{
			ItineraryName: "Sync Activities and Clear Themes",
			ThemeIds:      &emptyThemes,
			Days: &[]api.UpdateItineraryDayInput{
				{
					Id:         &day1.ID,
					DayNumber:  1,
					Title:      "Day 1 Updated",
					Subtitle:   nil,
					Activities: nil,
				},
				{
					Id:         &day2.ID,
					DayNumber:  2,
					Title:      "Day 2 Cleared Activities",
					Activities: &emptyActs,
				},
			},
		}
		req2 := shared.NewRequest(t, admin, http.MethodPut, fmt.Sprintf("/itineraries/%d", it.ID), payloadSync)
		w2 := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w2, req2)
		require.Equal(t, http.StatusOK, w2.Code, w2.Body.String())

		var themeCount2 int64
		shared.TestDB.Model(&dbmodel.ItineraryTheme{}).Where("itinerary_id = ?", it.ID).Count(&themeCount2)
		assert.Zero(t, themeCount2)

		var day1Acts, day2Acts []dbmodel.ItineraryActivity
		require.NoError(t, shared.TestDB.Where("itinerary_day_id = ?", day1.ID).Find(&day1Acts).Error)
		require.Len(t, day1Acts, 1)
		assert.Equal(t, act1.ID, day1Acts[0].ID)

		require.NoError(t, shared.TestDB.Where("itinerary_day_id = ?", day2.ID).Find(&day2Acts).Error)
		assert.Empty(t, day2Acts)
		var deletedAct2 dbmodel.ItineraryActivity
		assert.True(t, errors.Is(shared.TestDB.First(&deletedAct2, act2.ID).Error, gorm.ErrRecordNotFound))

		// 3. Empty Days slice: deletes all existing Days and Activities
		emptyDays := []api.UpdateItineraryDayInput{}
		payloadEmptyDays := api.UpdateItineraryRequest{
			ItineraryName: "Clear All Days",
			Days:          &emptyDays,
		}
		req3 := shared.NewRequest(t, admin, http.MethodPut, fmt.Sprintf("/itineraries/%d", it.ID), payloadEmptyDays)
		w3 := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w3, req3)
		require.Equal(t, http.StatusOK, w3.Code, w3.Body.String())

		var remainingDays []dbmodel.ItineraryDay
		require.NoError(t, shared.TestDB.Where("itinerary_id = ?", it.ID).Find(&remainingDays).Error)
		assert.Empty(t, remainingDays)

		var remainingActs []dbmodel.ItineraryActivity
		require.NoError(t, shared.TestDB.Find(&remainingActs).Error)
		assert.Empty(t, remainingActs)
	})

	t.Run("Error_NotFound", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		payload := api.UpdateItineraryRequest{
			ItineraryName: "Khám phá Đà Lạt",
		}
		req := shared.NewRequest(t, admin, http.MethodPut, "/itineraries/999999", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code, w.Body.String())
		assert.Contains(t, w.Body.String(), "Itinerary not found")
	})

	t.Run("Error_InvalidID", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		payload := api.UpdateItineraryRequest{
			ItineraryName: "Khám phá Đà Lạt",
		}
		req := shared.NewRequest(t, admin, http.MethodPut, "/itineraries/0", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
	})

	t.Run("Error_InvalidAndOtherItineraryDayID_And_DuplicateDayID", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		it1 := factory.NewItineraryFactory().WithName("Itinerary 1").MustInsert(shared.TestDB)
		it2 := factory.NewItineraryFactory().WithName("Itinerary 2").MustInsert(shared.TestDB)

		day1 := factory.NewItineraryDayFactory().WithItineraryID(it1.ID).WithDayNumber(1).WithTitle("Day 1").MustInsert(shared.TestDB)
		otherDay := factory.NewItineraryDayFactory().WithItineraryID(it2.ID).WithDayNumber(1).WithTitle("Other Day").MustInsert(shared.TestDB)

		// Other itinerary's day ID
		payloadOther := api.UpdateItineraryRequest{
			ItineraryName: "Itinerary 1",
			Days: &[]api.UpdateItineraryDayInput{
				{Id: &otherDay.ID, DayNumber: 1, Title: "Day 1"},
			},
		}
		req1 := shared.NewRequest(t, admin, http.MethodPut, fmt.Sprintf("/itineraries/%d", it1.ID), payloadOther)
		w1 := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusBadRequest, w1.Code, w1.Body.String())
		assert.Contains(t, w1.Body.String(), itineraryusecase.ErrUpdateItineraryInvalidDayID.Error())

		// Duplicate day ID
		payloadDupID := api.UpdateItineraryRequest{
			ItineraryName: "Itinerary 1",
			Days: &[]api.UpdateItineraryDayInput{
				{Id: &day1.ID, DayNumber: 1, Title: "Day 1"},
				{Id: &day1.ID, DayNumber: 2, Title: "Day 2"},
			},
		}
		req2 := shared.NewRequest(t, admin, http.MethodPut, fmt.Sprintf("/itineraries/%d", it1.ID), payloadDupID)
		w2 := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusBadRequest, w2.Code, w2.Body.String())
		assert.Contains(t, w2.Body.String(), itineraryusecase.ErrUpdateItineraryDuplicateDayID.Error())

		// Duplicate dayNumber
		payloadDupNum := api.UpdateItineraryRequest{
			ItineraryName: "Itinerary 1",
			Days: &[]api.UpdateItineraryDayInput{
				{Id: &day1.ID, DayNumber: 1, Title: "Day 1"},
				{DayNumber: 1, Title: "Day 1 Dup"},
			},
		}
		req3 := shared.NewRequest(t, admin, http.MethodPut, fmt.Sprintf("/itineraries/%d", it1.ID), payloadDupNum)
		w3 := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w3, req3)
		assert.Equal(t, http.StatusBadRequest, w3.Code, w3.Body.String())
		assert.Contains(t, w3.Body.String(), itineraryday.ErrDuplicateDayNumber.Error())
	})

	t.Run("Error_InvalidAndOtherDayActivityID_And_DuplicateActivityID", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		it := factory.NewItineraryFactory().WithName("Itinerary 1").MustInsert(shared.TestDB)
		day1 := factory.NewItineraryDayFactory().WithItineraryID(it.ID).WithDayNumber(1).WithTitle("Day 1").MustInsert(shared.TestDB)
		day2 := factory.NewItineraryDayFactory().WithItineraryID(it.ID).WithDayNumber(2).WithTitle("Day 2").MustInsert(shared.TestDB)

		act1 := factory.NewItineraryActivityFactory().WithItineraryDayID(day1.ID).WithStartTime("08:00:00").WithTitle("Act 1").MustInsert(shared.TestDB)
		act2 := factory.NewItineraryActivityFactory().WithItineraryDayID(day2.ID).WithStartTime("09:00:00").WithTitle("Act 2").MustInsert(shared.TestDB)

		// Activity belongs to day2 but sent under day1
		payloadOtherDayAct := api.UpdateItineraryRequest{
			ItineraryName: "Itinerary 1",
			Days: &[]api.UpdateItineraryDayInput{
				{
					Id:        &day1.ID,
					DayNumber: 1,
					Title:     "Day 1",
					Activities: &[]api.UpdateItineraryActivityInput{
						{Id: &act2.ID, StartTime: "08:00", Title: "Wrong Day Act"},
					},
				},
			},
		}
		req1 := shared.NewRequest(t, admin, http.MethodPut, fmt.Sprintf("/itineraries/%d", it.ID), payloadOtherDayAct)
		w1 := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusBadRequest, w1.Code, w1.Body.String())
		assert.Contains(t, w1.Body.String(), itineraryusecase.ErrUpdateItineraryInvalidActivityID.Error())

		// Duplicate activity ID
		payloadDupAct := api.UpdateItineraryRequest{
			ItineraryName: "Itinerary 1",
			Days: &[]api.UpdateItineraryDayInput{
				{
					Id:        &day1.ID,
					DayNumber: 1,
					Title:     "Day 1",
					Activities: &[]api.UpdateItineraryActivityInput{
						{Id: &act1.ID, StartTime: "08:00", Title: "Act 1"},
						{Id: &act1.ID, StartTime: "10:00", Title: "Act 1 Dup"},
					},
				},
			},
		}
		req2 := shared.NewRequest(t, admin, http.MethodPut, fmt.Sprintf("/itineraries/%d", it.ID), payloadDupAct)
		w2 := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusBadRequest, w2.Code, w2.Body.String())
		assert.Contains(t, w2.Body.String(), itineraryusecase.ErrUpdateItineraryDuplicateActivityID.Error())
	})

	t.Run("Error_DuplicateThemeID_And_InvalidImagePath_AndVerifyRollback", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		theme := factory.NewTravelThemeFactory().WithName("Thiên nhiên").MustInsert(shared.TestDB)
		it := factory.NewItineraryFactory().
			WithName("Original Name").
			WithCoverImagePath("itineraries/original-cover.jpg").
			MustInsert(shared.TestDB)

		// Duplicate theme ID
		dupThemes := []int64{theme.ID, theme.ID}
		payloadDupTheme := api.UpdateItineraryRequest{
			ItineraryName: "Updated Name",
			ThemeIds:      &dupThemes,
		}
		req1 := shared.NewRequest(t, admin, http.MethodPut, fmt.Sprintf("/itineraries/%d", it.ID), payloadDupTheme)
		w1 := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusBadRequest, w1.Code, w1.Body.String())
		assert.Contains(t, w1.Body.String(), itinerarytheme.ErrDuplicateThemeID.Error())

		// Invalid cover image path
		badCover := "homestays/bad-cover.jpg"
		payloadBadCover := api.UpdateItineraryRequest{
			ItineraryName:  "Updated Name",
			CoverImagePath: &badCover,
		}
		req2 := shared.NewRequest(t, admin, http.MethodPut, fmt.Sprintf("/itineraries/%d", it.ID), payloadBadCover)
		w2 := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusBadRequest, w2.Code, w2.Body.String())
		assert.Contains(t, w2.Body.String(), itinerary.ErrInvalidImagePath.Error())

		// Invalid activity image path causes transaction rollback
		badActImg := "itineraries/not-in-activities.jpg"
		payloadBadActImg := api.UpdateItineraryRequest{
			ItineraryName: "Should Rollback Name",
			Days: &[]api.UpdateItineraryDayInput{
				{
					DayNumber: 1,
					Title:     "Ngày 1",
					Activities: &[]api.UpdateItineraryActivityInput{
						{
							StartTime: "08:00",
							Title:     "Ăn sáng",
							ImagePath: &badActImg,
						},
					},
				},
			},
		}
		req3 := shared.NewRequest(t, admin, http.MethodPut, fmt.Sprintf("/itineraries/%d", it.ID), payloadBadActImg)
		w3 := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w3, req3)
		assert.Equal(t, http.StatusBadRequest, w3.Code, w3.Body.String())
		assert.Contains(t, w3.Body.String(), itineraryactivity.ErrInvalidImagePath.Error())

		var unchangedIt dbmodel.Itinerary
		require.NoError(t, shared.TestDB.First(&unchangedIt, it.ID).Error)
		assert.Equal(t, "Original Name", unchangedIt.ItineraryName)
	})

	t.Run("Error_InvalidBody_And_Unauthorized", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		reqBadBody := shared.NewRequest(t, admin, http.MethodPut, "/itineraries/1", strings.NewReader("{invalid-json"))
		w1 := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w1, reqBadBody)
		assert.Equal(t, http.StatusBadRequest, w1.Code)

		reqUnauth := shared.NewUnAuthRequest(t, http.MethodPut, "/itineraries/1", api.UpdateItineraryRequest{ItineraryName: "Test"})
		w2 := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w2, reqUnauth)
		assert.Equal(t, http.StatusUnauthorized, w2.Code)
	})
}
