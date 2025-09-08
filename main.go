package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	_ "shopflow/application/models"
	"shopflow/application/publisher"
	"shopflow/application/repository"
	"shopflow/application/routes"
	"shopflow/application/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"

	_ "shopflow/application/docs" // импорт сгенерированных swagger файлов

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Application service
// @version 1.0
// @description API приложения Application с JWT авторизацией
// @host localhost:8081
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

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

	// Подключение к Postgres через database/sql
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSL)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("[error] failed to connect to DB:", err)
	}
	defer db.Close()

	// Проверка соединения
	if err := db.Ping(); err != nil {
		log.Fatal("[error] cannot ping DB:", err)
	}

	// MQ eventPublisher
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatal("[error] failed to connect to RabbitMQ:", err)
	}
	defer conn.Close()
	appPublisher := publisher.NewApplicationPublisher(conn)

	// RabbitMQ
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	// DI
	appRepo := repository.NewApplicationRepository(db)
	appService := services.NewApplicationService(appRepo, appPublisher)
	eventPublisher := services.NewEventPublisher(conn)

	// Gin
	r := gin.Default()
	routes.RegisterApplicationRoutes(r, appService, eventPublisher)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Swagger с автоматическим сохранением токена
	swaggerURL := ginSwagger.URL("http://localhost:8081/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL, ginSwagger.PersistAuthorization(true)))

	log.Println("Application service running on port", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("[error] failed to run server:", err)
	}
}
