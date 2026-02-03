package usecases

import (
	"backend/models"
	"time"
)

type FeedingRepositoryScheduler interface {
	GetSchedulesByTime(timeStr string) ([]models.FeedingSchedule, error)
	DeductStockAndLogFeeding(items []ItemToFeed, dateTime time.Time, note string, createdBy uint) error
}
