package station

import (
	"encoding/json"
	"errors"
	"mrt-schedules/common/client"
	"net/http"
	"strings"
	"time"
)

type Service interface {
	GetAllStation() (response []StationResponse, err error)
	CheckScheduleByStation(id string) (response []ScheduleResponse, err error)
}

type service struct {
	client *http.Client
}

func NewService() Service {
	return &service{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *service) GetAllStation() (response []StationResponse, err error) {
	url := "https://jakartamrt.co.id/id/val/stasiuns"

	byteResponse, err := client.DoRequest(s.client, url)
	if err != nil {
		return
	}

	var stations []Station
	if err = json.Unmarshal(byteResponse, &stations); err != nil {
		return
	}

	for _, item := range stations {
		response = append(response, StationResponse{
			ID:   item.ID,
			Name: item.Name,
		})
	}

	return
}

func (s *service) CheckScheduleByStation(id string) (response []ScheduleResponse, err error) {
	url := "https://jakartamrt.co.id/id/val/stasiuns"

	byteResponse, err := client.DoRequest(s.client, url)
	if err != nil {
		return
	}

	var schedules []Schedule
	if err = json.Unmarshal(byteResponse, &schedules); err != nil {
		return
	}

	var selectedSchedule *Schedule
	for _, item := range schedules {
		if item.StationID == id {
			selectedSchedule = &item
			break
		}
	}

	if selectedSchedule == nil {
		err = errors.New("station not found")
		return
	}

	return ConvertDataToResponse(*selectedSchedule)
}

func ConvertDataToResponse(schedule Schedule) (response []ScheduleResponse, err error) {
	const (
		LebakBulusTripName = "Stasiun Lebak Bulus Grab"
		BundaranHITripName = "Stasiun Bundaran HI Bank DKI"
	)

	now := time.Now()

	// Parse schedule
	scheduleLebakBulusParsed, err := ConvertScheduleToTimeFormat(schedule.ScheduleLebakBulus)
	if err != nil {
		return
	}

	scheduleBundaranHIParsed, err := ConvertScheduleToTimeFormat(schedule.ScheduleBundaranHI)
	if err != nil {
		return
	}

	// Filter future schedule Lebak Bulus
	for _, item := range scheduleLebakBulusParsed {
		if item.After(now) {
			response = append(response, ScheduleResponse{
				StationName: LebakBulusTripName,
				Time:        item.Format("15:04:05"),
			})
		}
	}

	// Filter future schedule Bundaran HI
	for _, item := range scheduleBundaranHIParsed {
		if item.After(now) {
			response = append(response, ScheduleResponse{
				StationName: BundaranHITripName,
				Time:        item.Format("15:04:05"),
			})
		}
	}

	return
}

func ConvertScheduleToTimeFormat(schedule string) (response []time.Time, err error) {
	schedules := strings.Split(schedule, ",")
	now := time.Now()

	for _, item := range schedules {
		trimmedTime := strings.TrimSpace(item)
		if trimmedTime == "" {
			continue
		}

		// Parse dengan detik
		parsedTime, parseErr := time.Parse("15:04:05", trimmedTime)
		if parseErr != nil {
			err = errors.New("invalid time format: " + trimmedTime)
			return
		}

		// Gabungkan dengan tanggal hari ini
		fullDateTime := time.Date(
			now.Year(),
			now.Month(),
			now.Day(),
			parsedTime.Hour(),
			parsedTime.Minute(),
			parsedTime.Second(),
			0,
			now.Location(),
		)

		response = append(response, fullDateTime)
	}

	return
}
