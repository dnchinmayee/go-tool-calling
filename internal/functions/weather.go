package functions

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"tempfunctiontools/models"
)

func GetCurrentWeather(location string, format string) (string, error) {
	log.Println("Getting current weather for location", location, "in format", format)
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://wttr.in/"+location+"?format=%f", nil)
	log.Println("req", req)

	if err != nil {
		log.Printf("error creating request: %v", err)
		return "", err
	}
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("error calling API: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	var weather string

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response: %v", err)
		return "", err
	}
	weather = string(body)
	return fmt.Sprintf("Current weather for location %s in format %s is %s", location, format, weather), nil
}

func GetCurrentWeatherForeCast(location string, format string) (models.WeatherResponse, error) {
	log.Println("Getting current weather for location", location, "in format", format)
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://wttr.in/"+location+"?format=j1", nil)
	log.Println("req", req)

	result := models.WeatherResponse{}

	if err != nil {
		log.Printf("error creating request: %v", err)
		return result, err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error calling API: %v", err)
		return result, err
	}
	defer resp.Body.Close()

	// read the response into result
	var weather models.Weather
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		log.Printf("error decoding response: %v", err)
		return result, err
	}

	result = CreateResponse(format, weather)

	return result, nil

}

// create response based on format
func CreateResponse(format string, weather models.Weather) models.WeatherResponse {
	currentCondition := weather.CurrentCondition[0]

	resp := models.WeatherResponse{
		CurrentCondition: models.CurrentConditionResponse{
			FeelsLike:          currentCondition.FeelsLikeC,
			Temperature:        currentCondition.TempC,
			WeatherDescription: currentCondition.WeatherDesc[0].Value,
		},
	}

	if format == models.Fahrenheit {
		resp.CurrentCondition.FeelsLike = currentCondition.FeelsLikeF
		resp.CurrentCondition.Temperature = currentCondition.TempF
	}

	// loop through weather forecast and create response
	for _, forecast := range weather.WeatherForecast {
		forecastResponse := models.WeatherForecastResponse{
			Date:        forecast.Date,
			MaxTemp:     forecast.MaxTempC,
			MinTemp:     forecast.MinTempC,
			TotalSnowCm: forecast.TotalSnowCm,
			UvIndex:     forecast.UvIndex,
		}

		if format == models.Fahrenheit {
			forecastResponse.MaxTemp = forecast.MaxTempF
			forecastResponse.MinTemp = forecast.MinTempF
		}
		resp.WeatherForecast = append(resp.WeatherForecast, forecastResponse)
	}

	return resp
}
