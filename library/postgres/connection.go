package postgres

import (
	"campaign-service/models"
	"context"
	"database/sql"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var SQLDB *sql.DB

func ConnectDatabase(ctx context.Context, postgresConfig models.PostgresConfig) {

	// Construct the DSN string
	// dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s", postgresConfig.Host, postgresConfig.Port, postgresConfig.User, postgresConfig.Password, postgresConfig.DBName, postgresConfig.SSLMode, postgresConfig.TimeZone)

	dsn := postgresConfig.NeonDb

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

	if postgresConfig.NeonDb == "" {
		if postgresConfig.Host != "localhost" {

			// Create enum types first
			createEnumTypes(db)

			// CAMPAIGN TABLE
			db.AutoMigrate(&models.Campaign{})

			// LOCATION TABLE
			db.AutoMigrate(&models.Location{})

			// PARTICIPANT TABLE
			db.AutoMigrate(&models.Participant{})

			// CATEGORY TABLE
			db.AutoMigrate(&models.Category{})

			// CAMPAIGN INVITE TABLE
			db.AutoMigrate(&models.CampaignInvite{})

			// CAMPAIGN REVIEW TABLE
			db.AutoMigrate(&models.CampaignReview{})

			// STATUS LOGS TABLE
			db.AutoMigrate(&models.StatusLogs{})

		}
	}

	DB = db
	SQLDB = sqlDB

}

func createEnumTypes(db *gorm.DB) {
	db.Exec(`DO $$ BEGIN
		CREATE TYPE campaign_status AS ENUM ('draft', 'active', 'completed', 'inactive', 'cancelled', 'full');
	EXCEPTION
		WHEN duplicate_object THEN null;
	END $$;`)

	db.Exec(`DO $$ BEGIN
		CREATE TYPE participant_status AS ENUM ('pending', 'active', 'left', 'rejected');
	EXCEPTION
		WHEN duplicate_object THEN null;
	END $$;`)

	db.Exec(`DO $$ BEGIN
		CREATE TYPE status_type AS ENUM ('Active', 'Suspended', 'Deleted', 'Inactive', 'Pending Approval', 'Rejected', 'Approved', 'Submitted');
	EXCEPTION
		WHEN duplicate_object THEN null;
	END $$;`)

	log.Println("Enum types created successfully")
}
