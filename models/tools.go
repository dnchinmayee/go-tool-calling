package models

import (
	"fmt"
	"log"
	"strconv"
	"tempfunctiontools/internal/functions"
)

func GetWeatherTool(agent *Agent) Tool {
	return Tool{
		Type: "function",
		Function: &Function{
			Name:        "get_current_weather",
			Description: "Get the current weather for a location",
			Parameters: &Parameters{
				Type: "object",
				Properties: map[string]*Parameter{
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
			weather, err := functions.GetCurrentWeather(location, format)
			if err != nil {
				return nil, err
			}
			return map[string]any{
				"weather": weather,
			}, nil
		},
	}
}

func GetRevenueTool(agent *Agent) Tool {
	return Tool{
		Type: "function",
		Function: &Function{
			Name:        "get_revenue_by_month_and_year",
			Description: "Get the revenue by month and year",
			Parameters: &Parameters{
				Type: "object",
				Properties: map[string]*Parameter{
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
			// month := args["month"].(int)
			// year := args["year"].(int)
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
			revenue, err := functions.GetRevenue(month, year, agent.Db)
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

// func GetQuarterlyRevenueTool() Tool {
// 	return Tool{

// 		Type: "function",
// 		Function: &Function{
// 			Name:        "get_quarterly_revenue",
// 			Description: "Fetches total revenue for a given quarter and year",
// 			Parameters: &Parameters{
// 				Type: "object",
// 				Properties: map[string]*Parameter{
// 					"quarter": {Type: "integer", Description: "Quarter of the year (1-4)"},
// 					"year":    {Type: "integer", Description: "Year of revenue data"},
// 				},
// 				Required: []string{"quarter", "year"},
// 			},
// 		},
// 	}
// }

func GetTools(agent *Agent) []Tool {
	return []Tool{
		GetWeatherTool(agent),
		GetRevenueTool(agent),
	}
}

// register tools

func RegisterTools(agent *Agent) {
	for _, tool := range GetTools(agent) {
		agent.Tools[tool.Function.Name] = tool
	}
}
