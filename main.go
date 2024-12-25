// http api based on gin framework
// define all the http api

package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// http request handler
func httpHandler(c *gin.Context) {
	// get Parameter "topic" from http request
	topic := c.Query("topic")
	if topic == "" {
		topic = "default"
	}

	if !IsAKnownTopic(topic) {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "unknown topic"})
		return
	}
	triggerStartTopic := GetTriggerStartTopic()
	if triggerStartTopic == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "no trigger start topic"})
		return
	}
	value := 0.0
	switch topic {
	case triggerStartTopic:
		value = float64((time.Now().Unix() % 500) / 10) // the initial trigger
	case "TEST/MIRROR/CH01/TEMP":
		value = generateRandomNumber(400, 10)
	case "TEST/MIRROR/CH02/TEMP":
		value = generateRandomNumber(500, 5)
	case "TEST/MIRROR/CH03/TEMP":
		value = generateRandomNumber(600, 5)
	case "TEST/MIRROR/CH01/BELTSPEED":
		value = generateRandomNumber(0.2, 0.01)
	case "TEST/MIRROR/CH02/BELTSPEED":
		value = generateRandomNumber(0.3, 0.02)
	case "TEST/MIRROR/CH03/BELTSPEED":
		value = generateRandomNumber(0.4, 0.02)
	case "TEST/MIRROR/CH01/PRESSURE":
		value = generateRandomNumber(1.2, 0.02)
	default:
		value = generateRandomNumber(10, 0.5)
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok",
		"topic": topic, "value": value, "time": time.Now()})
}

func httpHandlerGetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "config": channels})
}

func httpHandlerSetConfig(c *gin.Context) {
	var newChannels Channels
	err := c.BindJSON(&newChannels)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}
	channels = newChannels
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func generateRandomNumber(level float64, noise float64) float64 {
	return level - noise/2 + noise*rand.Float64()
}

func main() {
	configPath := "/config/channels.json"
	flag.StringVar(&configPath, "config", configPath, "path to the channels.json file")
	flag.Parse()
	log.Println("Starting signal generator API server...")
	var err error
	log.Println("Reading channels.json file from ", configPath)
	channels, err = ReadChannels(configPath)
	if err != nil {
		log.Println("Error reading channels.json file: ", err.Error())
		os.Exit(1)
	} else {
		log.Println("Channels read successfully: ", len(channels.Tags), " variables and ", len(channels.Locations), " locations/triggers")
	}
	// create a default gin router
	router := gin.Default()

	// define the http api
	router.GET("/signal", httpHandler)
	router.POST("/config", httpHandlerSetConfig)
	router.GET("/config", httpHandlerGetConfig)
	port := os.Getenv("PORT_SIGNALSERVICE")
	if port == "" {
		port = "8091"
		log.Println("PORT_SIGNALSERVICE not set, using default port " + port)
	} else {
		log.Println("Port set to ", port)
	}
	// run the http server
	log.Println("starting signal generator at port ", port)
	router.Run(":" + port)
}
