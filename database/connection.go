package connection

import (
	"fmt"

	"github.com/aasoft24/golara/wpkg/configs"
	"github.com/aasoft24/golara/wpkg/database"
)

// Init uses already-loaded root YAML + root ENV for DB
func Init() error {
	// Root YAML already loaded in configs.GConfig
	dbCfg := configs.GConfig.Database
	if dbCfg.Default == "" {
		return fmt.Errorf("no default database connection found")
	}

	// Call wpkg database initializer
	if err := database.InitDB(); err != nil {
		return fmt.Errorf("database init failed: %v", err)
	}

	return nil
}

func GetDefaultConnection() string {
	return configs.GConfig.Database.Default
}
