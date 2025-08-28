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
	"golang.org/x/time/rate"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	db, err := initDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := db.AutoMigrate(&domain.User{}, &domain.Book{}, &domain.LendingRecord{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialise repositories
	userRepo := repository.NewUserRepository(db)
	bookRepo := repository.NewBookRepository(db)
	lendingRepo := repository.NewLendingRepository(db)

	// Initialise use cases
	authUC := usecase.NewAuthUseCase(userRepo)
	bookUC := usecase.NewBookUseCase(bookRepo)
	lendingUC := usecase.NewLendingUseCase(lendingRepo, bookRepo)

	// Initialise handlers
	jwtUtil := pkg.NewJWTUtil(cfg.JWT.Secret)
	authHandler := handler.NewAuthHandler(authUC, jwtUtil)
	bookHandler := handler.NewBookHandler(bookUC)
	lendingHandler := handler.NewLendingHandler(lendingUC)

	// Rate limiter – 100 requests per minute per IP with a burst of 200
	rateLimiter := middleware.NewRateLimiter(rate.Every(time.Minute/100), 200)

	// Setup router
	router := setupRouter(authHandler, bookHandler, lendingHandler, jwtUtil, rateLimiter)

	// Start server
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Server exited with error:", err)
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
	// attempt connection with retries
	for i := 0; i < 15; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			// Configure connection pool
			sqlDB, serr := db.DB()
			if serr != nil {
				return nil, serr
			}
			sqlDB.SetMaxIdleConns(10)
			sqlDB.SetMaxOpenConns(100)
			sqlDB.SetConnMaxLifetime(time.Hour)
			return db, nil
		}
		log.Println("Database not ready, retrying...")
		time.Sleep(2 * time.Second)
	}
	return nil, err
}

func setupRouter(
	authHandler *handler.AuthHandler,
	bookHandler *handler.BookHandler,
	lendingHandler *handler.LendingHandler,
	jwtUtil *pkg.JWTUtil,
	rateLimiter *middleware.RateLimiter,
) *gin.Engine {
	r := gin.Default()
	// global rate limiting
	r.Use(middleware.RateLimitMiddleware(rateLimiter))
	// simple CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	// health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "timestamp": time.Now()})
	})
	// API routes
	v1 := r.Group("/api/v1")
	// auth routes
	authGroup := v1.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
	}
	// book routes
	books := v1.Group("/books")
	{
		// listing and retrieval of books are public
		books.GET("", bookHandler.ListBooks)
		books.GET("/:id", bookHandler.GetBook)
		// modifications require authentication
		books.POST("", middleware.AuthMiddleware(jwtUtil), bookHandler.CreateBook)
		books.PUT("/:id", middleware.AuthMiddleware(jwtUtil), bookHandler.UpdateBook)
		books.DELETE("/:id", middleware.AuthMiddleware(jwtUtil), bookHandler.DeleteBook)
	}
	// lending routes
	lending := v1.Group("/lending").Use(middleware.AuthMiddleware(jwtUtil))
	{
		lending.POST("/borrow", lendingHandler.BorrowBook)
		lending.PUT("/return/:id", lendingHandler.ReturnBook)
		lending.GET("/history", lendingHandler.GetBorrowingHistory)
		lending.GET("/active", lendingHandler.GetActiveBorrowings)
	}
	return r
}
