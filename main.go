package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"orcaness.com/api/app"
)

func main() {
	viper.AddConfigPath("config/")
	viper.SetConfigName("app")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Read config file fail: %s", err)
	}

	defaultConfig := viper.AllSettings()

	envMap := map[string]string{
		"debug":   "dev",
		"release": "prod",
	}

	env := os.Getenv("GIN_MODE")
	if env == "" {
		env = os.Args[0]
	}

	if env != "release" {
		env = "debug"
	}

	gin.SetMode(env)

	// Load environment configration
	viper.SetConfigName(fmt.Sprintf("app.%s", envMap[env]))
	if err := viper.ReadInConfig(); err == nil {
		newConfig := viper.AllSettings()

		// Traverse up to 3 levels
		for key, value := range newConfig {
			if strings.HasPrefix(fmt.Sprintf("%T", value), "map") {
				for subKey, subValue := range value.(map[string]interface{}) {
					if defaultConfig[key] == nil {
						defaultConfig[key] = map[string]interface{}{}
					}

					if strings.HasPrefix(fmt.Sprintf("%T", subValue), "map") {
						for childSubKey, childSubValue := range subValue.(map[string]interface{}) {
							if defaultConfig[key].(map[string]interface{})[subKey] == nil {
								defaultConfig[key].(map[string]interface{})[subKey] = map[string]interface{}{}
							}
							defaultConfig[key].(map[string]interface{})[subKey].(map[string]interface{})[childSubKey] = childSubValue
						}
					} else {
						defaultConfig[key].(map[string]interface{})[subKey] = subValue
					}
				}
			} else {
				defaultConfig[key] = value
			}
		}

		// Reset viper configuration
		if configData, err := json.Marshal(defaultConfig); err == nil {
			viper.ReadConfig(bytes.NewBuffer(configData))
		}
	}

	router := gin.Default()

	router.Use(CORSMiddleware())
	app.CollectRoute(router)

	router.Run(":" + viper.GetString("server.port"))
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
