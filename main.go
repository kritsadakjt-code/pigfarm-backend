package main

import (
	"backend/adapters/handlers"
	"backend/adapters/repositories"
	"backend/adapters/schedulers"
	"backend/config"
	"backend/middlewares"
	"backend/models"
	"backend/routes"
	"backend/usecases"
	"backend/utils"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}
	config.InitDB()
	config.DB.AutoMigrate(
		&models.User{},
		&models.Pig{},
		&models.Breeding{},
		&models.FoodStock{},
		&models.Feeding{},
		&models.Health{},
		&models.Expense{},
		&models.PigSale{},
		&models.PigSaleItem{},
		&models.NotificationModel{},
		&models.FeedingSchedule{},
		&models.FeedingScheduleItem{},
	)
	// ฝัง owner
	utils.SeedOwnerUser(config.DB)

	// สร้าง Instance ของ FeedingService ใหม่
	feedingService := usecases.NewFeedingService(config.DB)

	notificationRepo := repositories.NewGormNotificationRepository(config.DB)
	foodStockRepo := repositories.NewGormFoodStockRepository(config.DB)
	breedingRepo := repositories.NewGormBreedingRepository(config.DB)

	notiService := usecases.NewNotificationService(
		notificationRepo,
		foodStockRepo,
		breedingRepo,
	)

	notiHandler := handlers.NewHttpNotificationHandler(notiService)

	// start scheduler
	scheduler := schedulers.NewNotificationScheduler(notiService, feedingService, 1*time.Minute)

	scheduler.Start()

	// graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		log.Println("\n shutting down gracefully...")
		scheduler.Stop()
		os.Exit(0)
	}()
	app := fiber.New()

	app.Use(middlewares.CorsConfig())
	routes.AuthRoutes(app)
	routes.ProfileRoutes(app)
	routes.EmailPasswordRoute(app)
	routes.UserByAdminRoute(app)
	routes.PigRoute(app)
	routes.BreedingRoute(app)
	routes.FeedingRoute(app)
	routes.FoodStockRoute(app)
	routes.HealthRoute(app)
	routes.ExpenseRoute(app)
	routes.PigSaleRoute(app)
	routes.DashboardRoute(app)
	routes.NotificationRoutes(app, notiHandler)
	routes.FeedingScheduleRoute(app)

	// utils.StartNotificationScheduler()
	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
