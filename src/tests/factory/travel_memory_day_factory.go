package factory

import (
	"time"

	"platform-go/internal/infrastructure/dbmodel"
)

// TravelMemoryDayFactory generates TravelMemoryDay dbmodel instances for seeders and tests.
type TravelMemoryDayFactory struct {
	memoryID    int64
	date        time.Time
	description *string
}

func NewTravelMemoryDayFactory() *TravelMemoryDayFactory {
	desc := "Ngày đầu tiên trong chuyến hành trình"
	return &TravelMemoryDayFactory{
		memoryID:    1,
		date:        time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
		description: &desc,
	}
}

func (f *TravelMemoryDayFactory) WithMemoryID(id int64) *TravelMemoryDayFactory {
	f.memoryID = id
	return f
}

func (f *TravelMemoryDayFactory) WithDate(date time.Time) *TravelMemoryDayFactory {
	f.date = date
	return f
}

func (f *TravelMemoryDayFactory) WithDescription(desc string) *TravelMemoryDayFactory {
	f.description = &desc
	return f
}

func (f *TravelMemoryDayFactory) Build() *dbmodel.TravelMemoryDay {
	return &dbmodel.TravelMemoryDay{
		MemoryID:    f.memoryID,
		Date:        f.date,
		Description: f.description,
	}
}

