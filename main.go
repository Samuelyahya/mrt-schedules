package main

import (
	"mrt-schedules/modules/station"

	"github.com/gin-gonic/gin"
)

func initiateRouter() {
	router := gin.Default()
	api := router.Group("/v1/api")

	station.Initiate(api)
	router.Run(":8000")
}

func main() {
	initiateRouter()
}
