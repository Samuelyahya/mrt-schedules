package station

import (
	"mrt-schedules/common/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Initiate(router *gin.RouterGroup) {
	stationService := NewService()
	station := router.Group("/station")
	station.GET("", func(c *gin.Context) {
		GetAllStation(c, stationService)
	})

	station.GET("/:id", func(c *gin.Context) {
		CheckScheduleByStation(c, stationService)
	})
}

func GetAllStation(c *gin.Context, service Service) {
	data, err := service.GetAllStation()
	if err != nil {
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})

		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Successfully get all station",
		Data:    data,
	})
}

func CheckScheduleByStation(c *gin.Context, service Service) {
	id := c.Param("id")
	data, err := service.CheckScheduleByStation(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Successfully get schedules by station",
		Data:    data,
	})
}
