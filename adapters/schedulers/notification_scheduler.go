package schedulers

import (
	"backend/config"
	"backend/models"
	"backend/usecases"
	"fmt"
	"log"
	"time"
)

type NotificationScheduler struct {
	service        *usecases.NotificationService
	feedingService *usecases.FeedingService
	interval       time.Duration
	ticker         *time.Ticker
	stopChan       chan struct{}
}

func NewNotificationScheduler(service *usecases.NotificationService, feedingService *usecases.FeedingService, interval time.Duration) *NotificationScheduler {
	return &NotificationScheduler{
		service:        service,
		feedingService: feedingService,
		interval:       interval,
		stopChan:       make(chan struct{}),
	}
}

func (s *NotificationScheduler) CheckFeedingSchedules() {
	log.Println("Checking for scheduled feeding tasks...")
	currentTime := time.Now().Format("15:04")

	var schedules []models.FeedingSchedule
	err := config.DB.Preload("Items").Where("is_active = ? AND scheduled_time = ?", true, currentTime).Find(&schedules).Error
	if err != nil {
		log.Printf("Error fetching feeding schedules: %v", err)
		return
	}

	if len(schedules) == 0 {
		return
	}

	for _, schedule := range schedules {
		if len(schedule.Items) == 0 {
			continue // ข้ามกฎที่ไม่มีรายการอาหาร
		}

		log.Printf("Executing schedule '%s' with %d item(s)", schedule.Name, len(schedule.Items))

		var itemsToFeed []usecases.ItemToFeed
		for _, item := range schedule.Items {
			itemsToFeed = append(itemsToFeed, usecases.ItemToFeed{
				FoodID: item.FoodID,
				Amount: item.Amount,
			})
		}

		err := s.feedingService.CreateFeedingLogsForItems(
			itemsToFeed,
			time.Now(),
			fmt.Sprintf("Automatic feeding : %s", schedule.Name), // สร้าง Note อัตโนมัติ
			schedule.CreatedBy,
		)

		if err != nil {
			log.Printf("ERROR executing schedule '%s': %v", schedule.Name, err)
			// (ทางเลือก) อาจจะสร้าง Notification แจ้งเตือน Admin ว่าการให้อาหารอัตโนมัติล้มเหลว
		}
	}
}

func (s *NotificationScheduler) Start() {
	s.ticker = time.NewTicker(s.interval)
	go func() {
		log.Printf("Notification scheduler started (run every %v)", s.interval)
		for {
			select {
			case <-s.ticker.C:
				log.Println("running scheduled notification checks...")
				if err := s.service.CheckFoodStock(); err != nil {
					log.Printf("error checking food stock: %v", err)
				}
				if err := s.service.CheckBreeding(); err != nil {
					log.Printf("error checking breeding: %v", err)
				}
				s.CheckFeedingSchedules()
			case <-s.stopChan:
				log.Println("Notification scheduler stopped")
				return
			}
		}
	}()
}

func (s *NotificationScheduler) Stop() {
	if s.ticker != nil {
		s.ticker.Stop()
		close(s.stopChan)
	}
}

// package schedulers

// import (
// 	"backend/config"
// 	"log"
// 	"time"
// )

// func StartNotificationScheduler() {
// 	ticker := time.NewTicker(1 * time.Minute)
// 	ns := adapters.repositories.NewNotificationService(config.DB)

// 	go func() {
// 		for range ticker.C {
// 			ns.RunAllChecks()
// 		}
// 	}()

// 	log.Println("Notification scheduler started (runs every hour)")
// }
