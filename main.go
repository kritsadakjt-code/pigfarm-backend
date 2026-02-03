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

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load() // โหลด .env
	config.InitDB()
	validate := validator.New()
	config.DB.AutoMigrate(
		&models.User{},
		&models.Pig{},
		&models.Breeding{},
		&models.FoodStock{},
		&models.FoodType{},
		&models.Feeding{},
		&models.FeedingItem{},
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

	notificationRepo := repositories.NewGormNotificationRepository(config.DB)
	foodStockRepo := repositories.NewGormFoodStockRepository(config.DB)
	breedingRepo := repositories.NewGormBreedingRepository(config.DB)
	feedingSchedRepo := repositories.NewFeedingScheduleGormRepo(config.DB)

	expenseRepo := repositories.NewExpenseGormRepository(config.DB)
	userRepo := repositories.NewUserGormRepository(config.DB)
	pigRepo := repositories.NewPigGormRepository(config.DB)
	breedingRepository := repositories.NewBreedingGormRepo(config.DB)
	feedingRepo := repositories.NewFeedingGormRepository(config.DB)
	stockRepo := repositories.NewStockGormRepo(config.DB)
	dashboardRepo := repositories.NewDashboardGormRepo(config.DB)

	notiService := usecases.NewNotificationService(
		notificationRepo,
		foodStockRepo,
		breedingRepo,
	)

	expenseService := usecases.NewExpenseService(expenseRepo)
	userService := usecases.NewUserService(userRepo)
	pigService := usecases.NewPigService(pigRepo)
	breedingService := usecases.NewBreedingService(breedingRepository, pigRepo)
	feedingService := usecases.NewFeedingService(feedingRepo)
	stockService := usecases.NewStockService(stockRepo)
	feedingSchedService := usecases.NewFeedingSchedulerService(feedingSchedRepo)
	dashboardService := usecases.NewDashboardService(dashboardRepo)

	notiHandler := handlers.NewHttpNotificationHandler(notiService)
	expenseHandler := handlers.NewHttpExpenseHandler(expenseService, validate)
	userHandler := handlers.NewHttpUserHandlers(userService, validate)
	pigHandler := handlers.NewHttpPigHandler(pigService, validate)
	breedingHandler := handlers.NewBreedingHttpHandler(breedingService, validate)
	feedingHandler := handlers.NewFeedingHttpHandler(feedingService)
	stockHandler := handlers.NewHttpStockHandler(stockService, validate)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)

	// start scheduler
	scheduler := schedulers.NewNotificationScheduler(notiService, feedingSchedService, 1*time.Minute)

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
	app := fiber.New(fiber.Config{
		ErrorHandler: middlewares.ErrorHandler,
	})

	app.Use(middlewares.CorsConfig())
	routes.AuthRoutes(app, userHandler)
	routes.ProfileRoutes(app)
	routes.EmailPasswordRoute(app, userHandler)
	routes.UserByAdminRoute(app, userHandler)
	routes.PigRoute(app, pigHandler)
	routes.BreedingRoute(app, breedingHandler)
	routes.FeedingRoute(app, feedingHandler)
	routes.FoodStockRoute(app, stockHandler)
	routes.FoodTypeRoute(app)
	routes.HealthRoute(app)
	routes.ExpenseRoute(app, expenseHandler)
	routes.PigSaleRoute(app)
	routes.DashboardRoute(app, dashboardHandler)
	routes.NotificationRoutes(app, notiHandler)
	routes.FeedingScheduleRoute(app)
	log.Println("FRONTEND_URL =", os.Getenv("FRONTEND_URL"))
	// utils.StartNotificationScheduler()
	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
