package functions

import (
	"fmt"
	"log"
	"strconv"
	"tempfunctiontools/models"
)

func GetWeatherTool(agent *models.Agent) models.Tool {
	return models.Tool{
		Type: "function",
		Function: &models.Function{
			Name:        "get_current_weather",
			Description: "Get the current weather for a location",
			Parameters: &models.Parameters{
				Type: "object",
				Properties: map[string]*models.Parameter{
					"location": {
						Type:        "string",
						Description: "The location to get the weather for, e.g. San Francisco, CA",
					},
					"format": {
						Type:        "string",
						Description: "The format to return the weather in, e.g. 'celsius' or 'fahrenheit'",
						Enum:        []string{"celsius", "fahrenheit"},
					},
				},
				Required: []string{"location", "format"},
			},
		},
		Execute: func(args map[string]any) (any, error) {
			location := args["location"].(string)
			format := args["format"].(string)
			weather, err := GetCurrentWeather(location, format)
			if err != nil {
				return nil, err
			}
			return map[string]any{
				"weather": weather,
			}, nil
		},
	}
}

func GetWeatherForecastTool(agent *models.Agent) models.Tool {
	return models.Tool{
		Type: "function",
		Function: &models.Function{
			Name:        "get_location_current_and_forecast_weather",
			Description: "Get a location's current and forecast weather",
			Parameters: &models.Parameters{
				Type: "object",
				Properties: map[string]*models.Parameter{
					"location": {
						Type:        "string",
						Description: "The location to get the weather for, e.g. San Francisco, CA",
					},
					"format": {
						Type:        "string",
						Description: "The format to return the weather in, e.g. 'celsius' or 'fahrenheit'",
						Enum:        []string{models.Celsius, models.Fahrenheit},
					},
				},
				Required: []string{"location", "format"},
			},
		},
		Execute: func(args map[string]any) (any, error) {
			location := args["location"].(string)
			format := args["format"].(string)
			weather, err := GetCurrentWeatherForeCast(location, format)
			if err != nil {
				return nil, err
			}
			return map[string]any{
				"location": location,
				"weather":  weather,
			}, nil
		},
	}
}

func GetRevenueTool(agent *models.Agent) models.Tool {
	return models.Tool{
		Type: "function",
		Function: &models.Function{
			Name:        "get_revenue_by_month_and_year",
			Description: "Get the revenue by month and year",
			Parameters: &models.Parameters{
				Type: "object",
				Properties: map[string]*models.Parameter{
					"month": {
						Type:        "integer",
						Description: "The month to get the revenue for, e.g. 1 for January, 2 for February, etc.",
					},
					"year": {
						Type:        "integer",
						Description: "The year to get the revenue for, e.g. 2023",
					},
				},
				Required: []string{"month", "year"},
			},
		},
		Execute: func(args map[string]any) (any, error) {
			month, err := strconv.Atoi(fmt.Sprintf("%v", args["month"]))
			if err != nil {
				log.Printf("error converting month to int: %v", err)
				return nil, err
			}

			year, err := strconv.Atoi(fmt.Sprintf("%v", args["year"]))
			if err != nil {
				log.Printf("error converting year to int: %v", err)
				return nil, err
			}
			revenue, err := GetRevenue(month, year, agent.Db)
			if err != nil {
				return nil, err
			}
			return map[string]any{
				"month":   month,
				"year":    year,
				"revenue": revenue,
			}, nil
		},
	}

}

func GetCurrentDateTimeLocationTool(agent *models.Agent) models.Tool {
	return models.Tool{
		Type: "function",
		Function: &models.Function{
			Name:        "get_current_location_date_time",
			Description: "Get the current location, date, time, and time zone information",
		},
		Execute: func(args map[string]any) (any, error) {
			dt := GetCurrentDateTimeLocation()
			return map[string]any{
				"date":     dt.Date,
				"time":     dt.Time,
				"location": dt.Location,
			}, nil
		},
	}
}

func GetTools(agent *models.Agent) []models.Tool {
	return []models.Tool{
		GetWeatherForecastTool(agent),
		GetRevenueTool(agent),
		GetCurrentDateTimeLocationTool(agent),
	}
}

// register tools

func RegisterTools(agent *models.Agent) {
	for _, tool := range GetTools(agent) {
		agent.Tools[tool.Function.Name] = tool
	}
}
