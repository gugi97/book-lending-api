package main

import (
	"book-lending-api/internal/config"
	"book-lending-api/internal/domain"
	"book-lending-api/internal/handler"
	"book-lending-api/internal/middleware"
	"book-lending-api/internal/repository"
	"book-lending-api/internal/usecase"
	"book-lending-api/pkg"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	db, err := initDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	if err = db.AutoMigrate(&domain.User{}, &domain.Book{}, &domain.LendingRecord{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	userRepo := repository.NewUserRepository(db)
	bookRepo := repository.NewBookRepository(db)
	lendingRepo := repository.NewLendingRepository(db)

	authUC := usecase.NewAuthUseCase(userRepo)
	bookUC := usecase.NewBookUseCase(bookRepo)
	lendingUC := usecase.NewLendingUseCase(lendingRepo, bookRepo)

	jwtUtil := pkg.NewJWTUtil(cfg.JWT.Secret)
	authHandler := handler.NewAuthHandler(authUC, jwtUtil)
	bookHandler := handler.NewBookHandler(bookUC)
	lendingHandler := handler.NewLendingHandler(lendingUC)

	rl := middleware.NewRateLimiter(rate.Every(time.Minute/100), 200)
	router := gin.Default()
	router.Use(middleware.RateLimitMiddleware(rl))
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "timestamp": time.Now()})
	})

	v1 := router.Group("/api/v1")
	authGroup := v1.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
	}
	books := v1.Group("/books")
	{
		books.GET("", bookHandler.ListBooks)
		books.GET("/:id", bookHandler.GetBook)
		books.POST("", middleware.AuthMiddleware(jwtUtil), bookHandler.CreateBook)
		books.PUT("/:id", middleware.AuthMiddleware(jwtUtil), bookHandler.UpdateBook)
		books.DELETE("/:id", middleware.AuthMiddleware(jwtUtil), bookHandler.DeleteBook)
	}
	lending := v1.Group("/lending").Use(middleware.AuthMiddleware(jwtUtil))
	{
		lending.POST("/borrow", lendingHandler.BorrowBook)
		lending.PUT("/return/:id", lendingHandler.ReturnBook)
		lending.GET("/history", lendingHandler.GetBorrowingHistory)
		lending.GET("/active", lendingHandler.GetActiveBorrowings)
	}

	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err = router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal(err)
	}
}

func initDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	var db *gorm.DB
	var err error
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			sqlDB, err := db.DB()
			if err != nil {
				return nil, err
			}
			sqlDB.SetMaxIdleConns(10)
			sqlDB.SetMaxOpenConns(100)
			sqlDB.SetConnMaxLifetime(time.Hour)
			return db, nil
		}
		time.Sleep(2 * time.Second)
	}
	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}
