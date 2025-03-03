package functions

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type DateTime struct {
	Date     string   `json:"date"`
	Time     string   `json:"time"`
	Location Location `json:"location,omitempty"`
}

type Location struct {
	Latitude  float32 `json:"latitude,omitempty"`
	Longitude float32 `json:"longitude,omitempty"`
	Country   string  `json:"country,omitempty"`
	City      string  `json:"city,omitempty"`
	Timezone  string  `json:"timezone,omitempty"`
}

type LocationData struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Latitude    float32 `json:"lat"`
	Longitude   float32 `json:"lon"`
	Timezone    string  `json:"timezone"`
}

// GetCurrentDateTime returns the current date and time
func GetCurrentDateTimeLocation() DateTime {
	log.Println("Getting current date and time")
	// load the current timezone
	dt := time.Now()

	// get the current date and time
	date := dt.Format("2006-01-02")
	time := dt.Format("15:04:05")
	timeZone := dt.Location().String()

	log.Println("timeZone", timeZone)

	result := DateTime{
		Date: date,
		Time: time,
	}

	loc, err := GetLocationInformation()
	if err != nil {
		log.Printf("error getting location: %v", err)
		return result
	}

	result.Location = loc

	// return the current date and time
	return result
}

// get location information
func GetLocationInformation() (Location, error) {
	log.Println("Getting location information")
	// get the current location
	loc, err := GetCurrentLocation()
	if err != nil {
		log.Printf("error getting location: %v", err)
		return Location{}, err
	}

	return Location{
		Latitude:  loc.Latitude,
		Longitude: loc.Longitude,
		Country:   loc.Country,
		City:      loc.City,
		Timezone:  loc.Timezone,
	}, nil
}

func GetCurrentLocation() (LocationData, error) {
	log.Println("Getting current location")
	// get the current location
	result := LocationData{}
	api := "http://ip-api.com/json"

	resp, err := http.Get(api)
	if err != nil {
		log.Printf("error calling API: %v", err)
		return result, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("error decoding response: %v", err)
		// print the response
		log.Printf("response: %+v", resp)
		return result, err
	}

	log.Printf("result: %+v", result)

	if result.Status != "success" {
		log.Printf("error getting location: %v", result.Status)
		return result, fmt.Errorf("error getting location: %v", result.Status)
	}

	return result, nil
}
