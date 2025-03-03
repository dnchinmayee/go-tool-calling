package functions

import (
	"fmt"
	"io"
	"log"
	"net/http"
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
