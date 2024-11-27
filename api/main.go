package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq" // postgres driver
	"github.com/vlegro/backend/api/controller"
	"github.com/vlegro/backend/api/repository"
	"github.com/vlegro/backend/api/service"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	servicePort := "3322"
	log.Printf("REST API started at %s...\n", servicePort)
	dbConnectionUrl, exists := os.LookupEnv("DB_CONNECTION_URL")
	if !exists {
		log.Fatal("DB_CONNECTION_URL env variable does not exist")
	}

	db, err := gorm.Open(postgres.Open(dbConnectionUrl))
	failOnError(err, "Could not open DB connection")
	dbConnection, err := db.DB()
	failOnError(err, "Could not get sql connection")

	defer dbConnection.Close()

	customerService := dependencyInjection(dbConnection)
	customerController := controller.NewCustomerController(customerService)
	err = http.ListenAndServe(fmt.Sprintf(":%s", servicePort), customerController.RestController())
	log.Fatal(err)
}

func dependencyInjection(dbConnection *sql.DB) *service.CustomerService {
	customerRepository := repository.NewCustomerRepositoryImpl(dbConnection)
	customerService := service.NewCustomerService(customerRepository)
	return customerService
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
