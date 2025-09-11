package main

import (
	_ "context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"shopflow/application/publisher"
	"shopflow/application/repository"
	"shopflow/application/routes"
	"shopflow/application/services"

	_ "shopflow/application/docs" // сгенерированные swagger файлы
	_ "shopflow/application/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Application Service
// @version 1.0
// @description API приложения Application с JWT авторизацией
// @host localhost:8081
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Загружаем переменные окружения
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// --- Подключение к Postgres ---
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_DATABASE")
	dbSSL := os.Getenv("DB_SSLMODE")
	if dbSSL == "" {
		dbSSL = "disable"
	}
	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		log.Fatal("[error] missing required DB environment variables")
	}
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSL)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("[error] failed to connect to DB:", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal("[error] cannot ping DB:", err)
	}

	// --- Подключение к RabbitMQ ---
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatal("[error] failed to connect to RabbitMQ:", err)
	}
	defer conn.Close()
	appPublisher := publisher.NewApplicationPublisher(conn)
	eventPublisher := services.NewEventPublisher(conn)

	// --- Подключение к gRPC Auth ---
	authGRPCAddr := os.Getenv("AUTH_GRPC_ADDR")
	if authGRPCAddr == "" {
		authGRPCAddr = "localhost:50051"
	}
	authClient, err := services.NewAuthClient(authGRPCAddr)
	if err != nil {
		log.Fatal("[error] failed to create AuthClient:", err)
	}

	// --- DI ---
	appRepo := repository.NewApplicationRepository(db)
	appService := services.NewApplicationService(appRepo, appPublisher, authClient)

	// --- Gin ---
	r := gin.Default()

	// Регистрируем маршруты приложения
	routes.RegisterApplicationRoutes(r, appService, eventPublisher)

	// Swagger
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	swaggerURL := ginSwagger.URL("http://localhost:" + port + "/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL, ginSwagger.PersistAuthorization(true)))

	log.Println("Application service running on port", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("[error] failed to run server:", err)
	}
}
