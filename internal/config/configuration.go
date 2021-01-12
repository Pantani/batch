package config

import (
	"strings"

	"github.com/Pantani/logger"
	_ "github.com/joho/godotenv/autoload" // auto import .env file variables
	"github.com/spf13/viper"
)

// configuration represents the env vars config object.
type configuration struct {
	Debug bool
	API   struct {
		Mode string
		Port string
	}
	Database struct {
		Type  string
		Redis struct {
			Host     string
			Password string
			Index    int
		}
	}
}

// Configuration configuration var to accessed globally.
var Configuration configuration

// print print all environment variables
func (c *configuration) print() {
	logger.Info("DEBUG", logger.Params{"value": c.Debug})
	logger.Info("API_PORT", logger.Params{"value": c.API.Port})
	logger.Info("API_MODE", logger.Params{"value": c.API.Mode})
	logger.Info("DATABASE_TYPE", logger.Params{"value": c.Database.Type})
	logger.Info("DATABASE_REDIS_HOST", logger.Params{"value": c.Database.Redis.Host})
	logger.Info("DATABASE_REDIS_PASSWORD", logger.Params{"value": c.Database.Redis.Password})
	logger.Info("DATABASE_REDIS_INDEX", logger.Params{"value": c.Database.Redis.Index})
}

// setDefaults set dummy values to force viper to search for these keys in environment variables
// the AutomaticEnv() only searches for already defined keys in a config file, default values or kvstore struct.
func setDefaults() {
	viper.SetDefault("Debug", true)
	viper.SetDefault("API.Port", "8889")
	viper.SetDefault("API.Mode", "debug")
	viper.SetDefault("Database.Type", "memory")
	viper.SetDefault("Database.Redis.Host", "localhost:6379")
	viper.SetDefault("Database.Redis.Password", "energix")
	viper.SetDefault("Database.Redis.Index", 1)
}

// Init reads in config file and ENV variables if set.
func Init() {
	setDefaults()
	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.Unmarshal(&Configuration); err != nil {
		logger.Error(err, "Error Unmarshal Viper Config File")
	}

	Configuration.print()
	logger.SetLogLevel(logger.InfoLevel)
	if Configuration.Debug {
		logger.SetLogLevel(logger.DebugLevel)
	}
}
