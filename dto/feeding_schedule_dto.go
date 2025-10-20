package dto

type FeedingScheduleItemInput struct {
	FoodID uint    `json:"food_id" validate:"required"`
	Amount float64 `json:"amount" `
}

type FeedingScheduleInput struct {
	Name          string                     `json:"name" validate:"required"`
	ScheduledTime string                     `json:"scheduled_time" validate:"required,datetime=15:04"`
	IsActive      bool                       `json:"is_active" `
	Note          string                     `json:"note"`
	Items         []FeedingScheduleItemInput `json:"items" validate:"required,min=1,dive"`
}
