package itinerary_controller_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	api "platform-go/generated/api"
	"platform-go/internal/controller_test/shared"
	"platform-go/internal/domain/itinerary"
	itineraryactivity "platform-go/internal/domain/itinerary_activity"
	itineraryday "platform-go/internal/domain/itinerary_day"
	itinerarytheme "platform-go/internal/domain/itinerary_theme"
	"platform-go/internal/infrastructure/dbmodel"
	"platform-go/internal/infrastructure/repository"
	itineraryusecase "platform-go/internal/usecases/itineraries/usecase"
	"platform-go/tests/factory"
)

func TestItineraryController_CreateItinerary(t *testing.T) {
	t.Run("Normal_CreateItinerarySuccess", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		loc := factory.NewTravelLocationFactory().
			WithLocationName("Lâm Đồng").
			WithCountryCode("VN").
			WithCountryName("Vietnam").
			MustInsert(shared.TestDB)

		theme1 := factory.NewTravelThemeFactory().
			WithName("Thiên nhiên").
			MustInsert(shared.TestDB)

		theme2 := factory.NewTravelThemeFactory().
			WithName("Ẩm thực").
			MustInsert(shared.TestDB)

		suitableFor := "Cặp đôi, nhóm bạn, gia đình"
		desc := "Một hành trình nhẹ nhàng để cảm nhận trọn vẹn vẻ đẹp của Đà Lạt."
		coverPath := "itineraries/dalat-cover.jpg"
		travelTips := []string{
			"Thời tiết Đà Lạt se lạnh vào sáng sớm và ban đêm, nên chuẩn bị áo ấm.",
			"Thuê xe máy sẽ giúp chủ động di chuyển.",
		}
		themeIDs := []int64{theme1.ID, theme2.ID}

		day1Subtitle := "Làm quen với thành phố ngàn hoa."
		day1Summary := "Cảm nhận không khí mát lành của trung tâm thành phố."
		day1Image := "itineraries/days/dalat-day-1.jpg"
		day1Tips := []string{"Nên mang theo áo khoác nhẹ."}

		actCategory := api.ItineraryActivityCategoryDining
		actDuration := 60
		actLocation := "Chợ Đà Lạt"
		actMaps := "https://www.google.com/maps?q=11.9427,108.4373"
		actDesc := "Thưởng thức bánh căn và sữa đậu nành nóng."
		actNotes := "Nên đi sớm trước 8h00."
		actImage := "itineraries/activities/dalat-banh-can.jpg"

		payload := api.CreateItineraryRequest{
			ItineraryName:  "Khám phá Đà Lạt",
			LocationId:     &loc.ID,
			SuitableFor:    &suitableFor,
			Description:    &desc,
			CoverImagePath: &coverPath,
			TravelTips:     &travelTips,
			ThemeIds:       &themeIDs,
			Days: &[]api.CreateItineraryDayInput{
				{
					DayNumber: 1,
					Title:     "Ngày 1 – Khám phá trung tâm Đà Lạt",
					Subtitle:  &day1Subtitle,
					Summary:   &day1Summary,
					ImagePath: &day1Image,
					DayTips:   &day1Tips,
					Activities: &[]api.CreateItineraryActivityInput{
						{
							StartTime:        "07:30",
							Title:            "Ăn sáng tại chợ Đà Lạt",
							ActivityCategory: &actCategory,
							DurationMinutes:  &actDuration,
							Location:         &actLocation,
							GoogleMapsUrl:    &actMaps,
							Description:      &actDesc,
							Notes:            &actNotes,
							ImagePath:        &actImage,
						},
					},
				},
				{
					DayNumber: 2,
					Title:     "Ngày 2 – Thiên nhiên và ngoại ô",
				},
			},
		}

		req := shared.NewRequest(t, admin, http.MethodPost, "/itineraries", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())

		var res api.MessageOnlyResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
		assert.Equal(t, itineraryusecase.CreateItineraryUseCaseSuccessMessage, res.Message)

		var it dbmodel.Itinerary
		require.NoError(t, shared.TestDB.Where("itinerary_name = ?", "Khám phá Đà Lạt").First(&it).Error)
		require.NotNil(t, it.LocationID)
		assert.Equal(t, loc.ID, *it.LocationID)
		require.NotNil(t, it.SuitableFor)
		assert.Equal(t, suitableFor, *it.SuitableFor)
		require.NotNil(t, it.Description)
		assert.Equal(t, desc, *it.Description)
		require.NotNil(t, it.CoverImagePath)
		assert.Equal(t, coverPath, *it.CoverImagePath)

		var themes []dbmodel.ItineraryTheme
		require.NoError(t, shared.TestDB.Where("itinerary_id = ?", it.ID).Order("id ASC").Find(&themes).Error)
		require.Len(t, themes, 2)
		assert.Equal(t, theme1.ID, themes[0].ThemeID)
		assert.Equal(t, theme2.ID, themes[1].ThemeID)

		var days []dbmodel.ItineraryDay
		require.NoError(t, shared.TestDB.Where("itinerary_id = ?", it.ID).Order("day_number ASC").Find(&days).Error)
		require.Len(t, days, 2)
		assert.Equal(t, 1, days[0].DayNumber)
		assert.Equal(t, "Ngày 1 – Khám phá trung tâm Đà Lạt", days[0].Title)
		assert.Equal(t, 2, days[1].DayNumber)

		var activities []dbmodel.ItineraryActivity
		require.NoError(t, shared.TestDB.Where("itinerary_day_id = ?", days[0].ID).Order("start_time ASC").Find(&activities).Error)
		require.Len(t, activities, 1)
		assert.Equal(t, "07:30:00", activities[0].StartTime)
		assert.Equal(t, "Ăn sáng tại chợ Đà Lạt", activities[0].Title)
		require.NotNil(t, activities[0].ActivityCategory)
		assert.Equal(t, "dining", *activities[0].ActivityCategory)
		require.NotNil(t, activities[0].DurationMinutes)
		assert.Equal(t, uint16(60), *activities[0].DurationMinutes)
	})

	t.Run("Normal_CreateItineraryWithoutOptionalFields", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		payload := api.CreateItineraryRequest{
			ItineraryName: "Lịch trình tự do",
		}

		req := shared.NewRequest(t, admin, http.MethodPost, "/itineraries", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())

		var it dbmodel.Itinerary
		require.NoError(t, shared.TestDB.Where("itinerary_name = ?", "Lịch trình tự do").First(&it).Error)
		assert.Nil(t, it.LocationID)
		assert.Nil(t, it.SuitableFor)
		assert.Nil(t, it.Description)
		assert.Nil(t, it.CoverImagePath)
		assert.Nil(t, it.TravelTips)
	})

	t.Run("Error_InvalidLocationID", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		invalidLocID := int64(999999)
		payload := api.CreateItineraryRequest{
			ItineraryName: "Khám phá Đà Lạt",
			LocationId:    &invalidLocID,
		}

		req := shared.NewRequest(t, admin, http.MethodPost, "/itineraries", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
		assert.Contains(t, w.Body.String(), "The specified travel location ID is invalid.")
	})

	t.Run("Error_InvalidThemeID", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		invalidThemeIDs := []int64{999999}
		payload := api.CreateItineraryRequest{
			ItineraryName: "Khám phá Đà Lạt",
			ThemeIds:      &invalidThemeIDs,
		}

		req := shared.NewRequest(t, admin, http.MethodPost, "/itineraries", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
		assert.Contains(t, w.Body.String(), "The specified travel theme ID is invalid.")
	})

	t.Run("Error_DuplicateThemeID", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		theme := factory.NewTravelThemeFactory().
			WithName("Thiên nhiên").
			MustInsert(shared.TestDB)

		dupThemes := []int64{theme.ID, theme.ID}
		payload := api.CreateItineraryRequest{
			ItineraryName: "Khám phá Đà Lạt",
			ThemeIds:      &dupThemes,
		}

		req := shared.NewRequest(t, admin, http.MethodPost, "/itineraries", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
		assert.Contains(t, w.Body.String(), itinerarytheme.ErrDuplicateThemeID.Error())
	})

	t.Run("Error_DuplicateDayNumber", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		payload := api.CreateItineraryRequest{
			ItineraryName: "Khám phá Đà Lạt",
			Days: &[]api.CreateItineraryDayInput{
				{DayNumber: 1, Title: "Ngày 1"},
				{DayNumber: 1, Title: "Ngày 1 trùng"},
			},
		}

		req := shared.NewRequest(t, admin, http.MethodPost, "/itineraries", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
		assert.Contains(t, w.Body.String(), itineraryday.ErrDuplicateDayNumber.Error())
	})

	t.Run("Error_InvalidStartTime", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		payload := api.CreateItineraryRequest{
			ItineraryName: "Khám phá Đà Lạt",
			Days: &[]api.CreateItineraryDayInput{
				{
					DayNumber: 1,
					Title:     "Ngày 1",
					Activities: &[]api.CreateItineraryActivityInput{
						{
							StartTime: "25:99",
							Title:     "Hoạt động lỗi giờ",
						},
					},
				},
			},
		}

		req := shared.NewRequest(t, admin, http.MethodPost, "/itineraries", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
		assert.Contains(t, w.Body.String(), itineraryactivity.ErrInvalidStartTime.Error())
	})

	t.Run("Error_InvalidCoverImagePath", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		badCover := "homestays/550e8400-e29b-41d4-a716-446655440000.jpg"
		payload := api.CreateItineraryRequest{
			ItineraryName:  "Khám phá Đà Lạt",
			CoverImagePath: &badCover,
		}

		req := shared.NewRequest(t, admin, http.MethodPost, "/itineraries", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
		assert.Contains(t, w.Body.String(), itinerary.ErrInvalidImagePath.Error())
	})

	t.Run("Error_InvalidDayImagePath", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		badDayImg := "itineraries/day-without-days-folder.jpg"
		payload := api.CreateItineraryRequest{
			ItineraryName: "Khám phá Đà Lạt",
			Days: &[]api.CreateItineraryDayInput{
				{
					DayNumber: 1,
					Title:     "Ngày 1",
					ImagePath: &badDayImg,
				},
			},
		}

		req := shared.NewRequest(t, admin, http.MethodPost, "/itineraries", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
		assert.Contains(t, w.Body.String(), itineraryday.ErrInvalidImagePath.Error())
	})

	t.Run("Error_InvalidActivityImagePath", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		badActImg := "itineraries/act-without-activities-folder.jpg"
		payload := api.CreateItineraryRequest{
			ItineraryName: "Khám phá Đà Lạt",
			Days: &[]api.CreateItineraryDayInput{
				{
					DayNumber: 1,
					Title:     "Ngày 1",
					Activities: &[]api.CreateItineraryActivityInput{
						{
							StartTime: "08:00",
							Title:     "Ăn sáng",
							ImagePath: &badActImg,
						},
					},
				},
			},
		}

		req := shared.NewRequest(t, admin, http.MethodPost, "/itineraries", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())
		assert.Contains(t, w.Body.String(), itineraryactivity.ErrInvalidImagePath.Error())
	})

	t.Run("Error_InvalidBody", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		req := shared.NewRequest(t, admin, http.MethodPost, "/itineraries", strings.NewReader("{invalid-json"))
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error_Unauthorized", func(t *testing.T) {
		shared.SetupTestEnvironment(t)

		payload := api.CreateItineraryRequest{
			ItineraryName: "Khám phá Đà Lạt",
		}

		req := shared.NewUnAuthRequest(t, http.MethodPost, "/itineraries", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Normal_AssociateActivitiesRegardlessOfDayOrder", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		payload := api.CreateItineraryRequest{
			ItineraryName: "Khám phá Đà Lạt ngược thứ tự ngày",
			Days: &[]api.CreateItineraryDayInput{
				{
					DayNumber: 2,
					Title:     "Ngày 2",
					Activities: &[]api.CreateItineraryActivityInput{
						{StartTime: "09:00", Title: "Hoạt động Ngày 2"},
					},
				},
				{
					DayNumber: 1,
					Title:     "Ngày 1",
					Activities: &[]api.CreateItineraryActivityInput{
						{StartTime: "07:30", Title: "Hoạt động Ngày 1"},
					},
				},
			},
		}

		req := shared.NewRequest(t, admin, http.MethodPost, "/itineraries", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code, w.Body.String())

		var it dbmodel.Itinerary
		require.NoError(t, shared.TestDB.Where("itinerary_name = ?", "Khám phá Đà Lạt ngược thứ tự ngày").First(&it).Error)

		var day1, day2 dbmodel.ItineraryDay
		require.NoError(t, shared.TestDB.Where("itinerary_id = ? AND day_number = 1", it.ID).First(&day1).Error)
		require.NoError(t, shared.TestDB.Where("itinerary_id = ? AND day_number = 2", it.ID).First(&day2).Error)

		var act1, act2 dbmodel.ItineraryActivity
		require.NoError(t, shared.TestDB.Where("itinerary_day_id = ?", day1.ID).First(&act1).Error)
		assert.Equal(t, "Hoạt động Ngày 1", act1.Title)

		require.NoError(t, shared.TestDB.Where("itinerary_day_id = ?", day2.ID).First(&act2).Error)
		assert.Equal(t, "Hoạt động Ngày 2", act2.Title)
	})

	t.Run("Error_TransactionRollbackOnFailure", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		admin := shared.CreateLoginUser()

		theme := factory.NewTravelThemeFactory().
			WithName("Thiên nhiên").
			MustInsert(shared.TestDB)

		themeIDs := []int64{theme.ID}
		payload := api.CreateItineraryRequest{
			ItineraryName: "Lịch trình lỗi rollback",
			ThemeIds:      &themeIDs,
			Days: &[]api.CreateItineraryDayInput{
				{
					DayNumber: 1,
					Title:     "Ngày 1 hợp lệ",
					Activities: &[]api.CreateItineraryActivityInput{
						{StartTime: "29:99", Title: "Hoạt động lỗi giờ gây rollback"},
					},
				},
			},
		}

		req := shared.NewRequest(t, admin, http.MethodPost, "/itineraries", payload)
		w := httptest.NewRecorder()
		shared.TestRouter.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, w.Body.String())

		var itCount, themeCount, dayCount, actCount int64
		shared.TestDB.Model(&dbmodel.Itinerary{}).Count(&itCount)
		shared.TestDB.Model(&dbmodel.ItineraryTheme{}).Count(&themeCount)
		shared.TestDB.Model(&dbmodel.ItineraryDay{}).Count(&dayCount)
		shared.TestDB.Model(&dbmodel.ItineraryActivity{}).Count(&actCount)

		assert.Zero(t, itCount)
		assert.Zero(t, themeCount)
		assert.Zero(t, dayCount)
		assert.Zero(t, actCount)
	})
}

func TestItineraryRepositories(t *testing.T) {
	t.Run("Repository_BulkEmptySlices_UpdateZeroValues_InvalidJSON", func(t *testing.T) {
		shared.SetupTestEnvironment(t)
		ctx := context.Background()
		schemaDB := repository.NewSchemaDB(shared.TestDB)

		itRepo := repository.NewItineraryRepository(schemaDB)
		thRepo := repository.NewItineraryThemeRepository(schemaDB)
		dayRepo := repository.NewItineraryDayRepository(schemaDB)
		actRepo := repository.NewItineraryActivityRepository(schemaDB)

		// Bulk insert & delete with empty slices should succeed as no-ops
		require.NoError(t, thRepo.CreateItineraryThemes(ctx, nil))
		require.NoError(t, dayRepo.CreateItineraryDays(ctx, nil))
		require.NoError(t, actRepo.CreateItineraryActivities(ctx, nil))
		require.NoError(t, dayRepo.DeleteDaysByIDs(ctx, nil))
		require.NoError(t, actRepo.DeleteActivitiesByDayIDs(ctx, nil))
		require.NoError(t, actRepo.DeleteActivitiesByIDs(ctx, nil))

		count, err := thRepo.CountValidThemesByIDs(ctx, nil)
		require.NoError(t, err)
		assert.Zero(t, count)

		// Create itinerary with populated tips, then update with nil/zero-value optional fields
		desc := "Mô tả ban đầu"
		itEntity, err := itinerary.NewItinerary("Tour Đà Lạt", nil, nil, &desc, nil, []string{"Tip 1"})
		require.NoError(t, err)
		require.NoError(t, itRepo.CreateItinerary(ctx, &itEntity))

		found, err := itRepo.FindItineraryByID(ctx, itEntity.ID)
		require.NoError(t, err)
		require.NotNil(t, found)
		assert.Equal(t, []string{"Tip 1"}, found.TravelTips)
		require.NotNil(t, found.Description)

		// Update clearing optional fields to nil
		require.NoError(t, itEntity.Update("Tour Đà Lạt Updated", nil, nil, nil, nil, nil))
		require.NoError(t, itRepo.UpdateItinerary(ctx, &itEntity))

		foundAfter, err := itRepo.FindItineraryByID(ctx, itEntity.ID)
		require.NoError(t, err)
		require.NotNil(t, foundAfter)
		assert.Equal(t, "Tour Đà Lạt Updated", foundAfter.ItineraryName)
		assert.Nil(t, foundAfter.Description)
		assert.Nil(t, foundAfter.TravelTips)

		// Updating a non-existent itinerary returns ErrItineraryNotFound
		missingIt := itinerary.Itinerary{ID: 999999, ItineraryName: "Missing"}
		require.ErrorIs(t, itRepo.UpdateItinerary(ctx, &missingIt), itinerary.ErrItineraryNotFound)

		// Updating a non-existent day returns ErrDayNotFound
		missingDay := itineraryday.ItineraryDay{ID: 999999, ItineraryID: itEntity.ID, DayNumber: 1, Title: "Missing"}
		require.ErrorIs(t, dayRepo.UpdateItineraryDays(ctx, []itineraryday.ItineraryDay{missingDay}), itineraryday.ErrDayNotFound)

		// Updating a non-existent activity returns ErrActivityNotFound
		missingAct := itineraryactivity.ItineraryActivity{ID: 999999, ItineraryDayID: 999999, StartTime: "08:00:00", Title: "Missing"}
		require.ErrorIs(t, actRepo.UpdateItineraryActivities(ctx, []itineraryactivity.ItineraryActivity{missingAct}), itineraryactivity.ErrActivityNotFound)
	})
}
