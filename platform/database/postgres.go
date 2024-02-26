package database

import (
	"fmt"
	"idealist/pkg/utils"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// MysqlConnection func for connection to Mysql database.
func PostgresSQLConnection() (*gorm.DB, error) {
	// Define database connection settings.
	maxConn, _ := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS"))
	maxIdleConn, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNECTIONS"))
	maxLifetimeConn, _ := strconv.Atoi(os.Getenv("DB_MAX_LIFETIME_CONNECTIONS"))

	// Build Mysql connection URL.
	postgresConnURL, err := utils.ConnectionURLBuilder("postgres")
	if err != nil {
		return nil, err
	}

	// Define database connection for MySQL through GORM.
	db, err := gorm.Open(postgres.Open(postgresConnURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, fmt.Errorf("error, not connected to database, %w", err)
	}

	// Set logger for future operations.
	db.Logger = db.Logger.LogMode(logger.Info)

	// Set database connection settings:
	// 	- SetMaxOpenConns: the default is 0 (unlimited)
	// 	- SetMaxIdleConns: defaultMaxIdleConns = 2
	// 	- SetConnMaxLifetime: 0, connections are reused forever
	psqlDB, err := db.DB()

	if err != nil {
		return nil, fmt.Errorf("error, not connected to database, %w", err)
	}

	psqlDB.SetMaxOpenConns(maxConn)
	psqlDB.SetMaxIdleConns(maxIdleConn)
	psqlDB.SetConnMaxLifetime(time.Duration(maxLifetimeConn))

	// Try to ping database.
	if err := psqlDB.Ping(); err != nil {
		defer psqlDB.Close() // close database connection
		return nil, fmt.Errorf("error, not sent ping to database, %w", err)
	}

	return db, nil
}
