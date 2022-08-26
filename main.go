package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/drzombey/aur-package-builder-api/database"
	"github.com/drzombey/aur-package-builder-api/handler"
	"github.com/drzombey/aur-package-builder-api/model"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	app model.App
)

func main() {
	setupLogFormatter()
	loadConfig()
	database.SetupMongoDB(&app)
	setupWebserver()
}

func setupWebserver() {
	gin.SetMode(gin.ReleaseMode)

	if app.Config.Debug {
		gin.SetMode(gin.DebugMode)
	}

	server := gin.Default()
	registerHandlers(server)
	server.Run(fmt.Sprintf(":%d", app.Config.WebserverPort))
}

func setupLogFormatter() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
}

func loadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	viper.SetDefault("webserverPort", 8080)
	viper.SetDefault("webserverMode", "production")
	viper.SetDefault("database", map[string]interface{}{
		"host":         "localhost",
		"port":         27017,
		"user":         "root",
		"password":     "example",
		"databaseName": "aur_package_builder",
	})

	if err := viper.ReadInConfig(); err != nil {
		log.Warnf("Error reading config file, using default values. %s", err)
	}

	err := viper.Unmarshal(&app.Config)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	logLevel, err := log.ParseLevel(app.Config.LogLevel)
	if err == nil {
		log.SetLevel(logLevel)
	}

	log.Info("Config loaded.")
}

func registerHandlers(s *gin.Engine) {
	s.GET("/package", handler.HandleGetPackage)
}
