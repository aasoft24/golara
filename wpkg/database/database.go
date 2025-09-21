// wpkg/database/database.go
package database

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/aasoft24/golara/wpkg/configs"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *gorm.DB
var MongoClient *mongo.Client

func InitDB() error {
	cfg := configs.GConfig
	dbConn := cfg.Database.Default
	conn := cfg.Database.Connections[dbConn]

	var err error

	// Get timezone from config
	timezone := cfg.App.Timezone
	if timezone == "" {
		timezone = "Local"
	}

	switch dbConn {
	case "mysql":
		loc := url.QueryEscape(timezone)
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=%s",
			conn["username"], conn["password"], conn["host"], conn["port"], conn["database"], loc)
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	case "pgsql":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
			conn["host"], conn["username"], conn["password"], conn["database"], conn["port"], timezone)
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	case "sqlite":
		DB, err = gorm.Open(sqlite.Open(conn["path"]), &gorm.Config{})

	case "sqlserver":
		dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
			conn["username"], conn["password"], conn["host"], conn["port"], conn["database"])
		DB, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})

	case "mongodb":
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		clientOptions := options.Client().ApplyURI(conn["uri"])
		MongoClient, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			return errors.New("❌ MongoDB connect failed: " + err.Error())
		}
		err = MongoClient.Ping(ctx, nil)
		if err != nil {
			return errors.New("❌ MongoDB ping failed: " + err.Error())
		}
		fmt.Println("✅ MongoDB connected successfully")
		return nil

	default:
		return errors.New("❌ Invalid DB connection type: " + dbConn)
	}

	if err != nil {
		return errors.New("❌ Failed to connect " + dbConn + ": " + err.Error())
	}

	fmt.Println("✅ " + dbConn + " connected successfully")
	return nil
}
