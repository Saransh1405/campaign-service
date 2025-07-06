package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"users-service/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var SQLDB *sql.DB

func ConnectDatabase(ctx context.Context, postgresConfig models.PostgresConfig) {

	// Construct the DSN string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s", postgresConfig.Host, postgresConfig.Port, postgresConfig.User, postgresConfig.Password, postgresConfig.DBName, postgresConfig.SSLMode, postgresConfig.TimeZone)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Could not connect!", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Could not connect to PostgresDB:\n", err)
	}
	err = sqlDB.PingContext(ctx)
	if err != nil {
		log.Fatal("Could not connect to PostgresDB:\n", err)
	} else {
		log.Println("Connected to Postgres")
	}

	//migration of Tables

	if postgresConfig.Host != "localhost" {

		// // LOGS TABLE
		// db.AutoMigrate(&models.Logs{})

		// // STATUS LOGS TABLE
		// db.AutoMigrate(&models.StatusLogs{})

		// // ROLE TABLE
		// db.AutoMigrate(&models.Role{})

		// // ACCOUNTS DATA TABLE
		// db.AutoMigrate(&models.AccountsData{})

		// // USER TABLE
		// db.AutoMigrate(&models.Users{})

		// // SUPPORTED LANGUAGES TABLE
		// db.AutoMigrate(&models.SupportedLanguages{})

		// // CURRENCIES TABLE
		// db.AutoMigrate(&models.Currencies{})

		// // STATE TAXES TABLE
		// db.AutoMigrate(&models.StateTaxes{})

		// // COUNTRIES TABLE
		// db.AutoMigrate(&models.Countries{})

		// // COUNTRY LEVEL TAXES TABLE
		// db.AutoMigrate(&models.CountryLevelTaxes{})

		// // TAX FIELDS TABLE
		// db.AutoMigrate(&models.TaxFields{})

		// // STATES TABLE
		// db.AutoMigrate(&models.States{})

		// // BANK ACCOUNT FIELDS TABLE
		// db.AutoMigrate(&models.BankAccountFields{})

		// // PROPERTY AMENITIES TABLE
		// db.AutoMigrate(&models.PropertyAmenities{})

		// // ROOM VIEWS TABLE
		// db.AutoMigrate(&models.RoomViews{})

		// // PROPERTY TYPE TABLE
		// db.AutoMigrate(&models.PropertyType{})

		// // FIELDS TABLE
		// db.AutoMigrate(&models.Fields{})

		// // BUSINESS TABLE
		// db.AutoMigrate(&models.Businesses{})

		// // ADDRESS TABLE
		// db.AutoMigrate(&models.Address{})

		// // BRANDS TABLE
		// db.AutoMigrate(&models.Brands{})

		// // PROPERTIES TABLE
		// db.AutoMigrate(&models.Properties{})

		// // PHOTOS TABLE
		// db.AutoMigrate(&models.Media{})

		// // CONTACT PERSON TABLE
		// db.AutoMigrate(&models.ContactPerson{})

		// // // ACCOUNTS TABLE
		// // db.AutoMigrate(&models.Accounts{})

		// // DOCUMENTS FOR ENLISTING TABLE
		// db.AutoMigrate(&models.DocumentsForEnlisting{})

		// // // WORKING HOURS TABLE
		// // db.AutoMigrate(&models.WorkingHours{})

		// // // HOLIDAYS TABLE
		// // db.AutoMigrate(&models.Holidays{})

		// // // LEAVE TYPES TABLE
		// // db.AutoMigrate(&models.LeaveTypes{})

		// // REASONS TABLE
		// db.AutoMigrate(&models.Reasons{})

		// // BUSINEES META DATA TABLE
		// db.AutoMigrate(&models.BusinessMetaData{})

	}

	DB = db
	SQLDB = sqlDB

}
