package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"tempfunctiontools/internal/database"
	"tempfunctiontools/models"

	"github.com/gin-gonic/gin"
)

type ChatController struct {
	ctx    context.Context
	db     *database.DbConfig
	agent  *models.Agent
	apiKey string
}

func NewChatController(ctx context.Context, agent *models.Agent, db *database.DbConfig) *ChatController {
	agent.Db = db
	models.RegisterTools(agent)

	return &ChatController{
		ctx:   ctx,
		db:    db,
		agent: agent,
	}
}

func (ctrl *ChatController) GetChat(c *gin.Context) {
	apiKey := c.Request.Header.Get("Authorization")
	apiKey = apiKey[7:]

	chatBody := models.ChatBody{}

	if err := c.ShouldBindJSON(&chatBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println(apiKey)
	log.Println(chatBody)

	ctrl.apiKey = apiKey

	returnMessages, err := ctrl.ProcessQuery(c, chatBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, returnMessages)
}

func (ctrl *ChatController) GetChat1(c *gin.Context, agent *models.Agent) {
	apiKey := c.Request.Header.Get("Authorization")
	apiKey = apiKey[7:]

	chatBody := models.ChatBody{}

	if err := c.ShouldBindJSON(&chatBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println(apiKey)
	log.Println(chatBody)

	chatBody.Tools = models.GetTools(ctrl.agent)

	// chatBody.Messages = append([]models.Message{
	// 	{
	// 		Role:    "system",
	// 		Content: "You have access to the following tools: get_revenue_by_month_and_year. Always use these tools when applicable.",
	// 	},
	// }, chatBody.Messages...)

	for _, toolCall := range chatBody.Tools {
		if toolCall.Function.Name == "get_quarterly_revenue" {

			if toolCall.Function.Description == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Function arguments are missing"})
				return
			}
			// Ensure parameters are parsed correctly
			var args map[string]int
			jsonBytes, err := json.Marshal(toolCall.Function.Description)

			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid function arguments"})
				return
			}
			if err := json.Unmarshal(jsonBytes, &args); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid function arguments"})
				return
			}

			quarter, quarterOk := args["quarter"]
			year, yearOk := args["year"]
			if !quarterOk || !yearOk {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Missing quarter or year argument"})
				return
			}

			// Get revenue for the requested quarter
			revenue, err := ctrl.GetQuarterlyRevenueInternal(quarter, year)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Return the revenue response
			c.JSON(http.StatusOK, gin.H{
				"tool":    "get_quarterly_revenue",
				"quarter": quarter,
				"year":    year,
				"revenue": revenue,
			})
			return
		}
	}

	// Otherwise, forward request to LLM
	response, err := ctrl.callLLM(ctrl.ctx, chatBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

func (ctrl *ChatController) callLLM(ctx context.Context, chatBody models.ChatBody) (models.ChatResponse, error) {
	responseBody := models.ChatResponse{}

	client := &http.Client{}

	log.Printf("chatBody: %+v", chatBody)

	jsonBytes, err := json.Marshal(chatBody)
	if err != nil {
		log.Printf("error marshalling json: %v", err)
		return responseBody, err
	}

	log.Printf("jsonBytes: %s", jsonBytes)

	// req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonBytes))
	req, err := http.NewRequestWithContext(ctx, "POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+ctrl.apiKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error calling LLM: %v", err)
		return responseBody, err
	}

	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		log.Printf("error decoding response: %v", err)
		return responseBody, err
	}

	log.Printf("responseBody: %+v", responseBody)

	return responseBody, nil
}

// func (ctrl *ChatController) extractFunctionCall(toolCalls []models.ToolCall) []models.Message {

// 	var messages []models.Message

// 	for _, toolCall := range toolCalls {
// 		functionName := toolCall.Function.Name
// 		arguments := toolCall.Function.Arguments

// 		log.Printf("Function name: %s, Arguments: %v", functionName, arguments)

// 		// Ensure parameters are parsed correctly
// 		jsonBytes, err := json.Marshal(arguments)
// 		if err != nil {
// 			log.Printf("error marshalling arguments: %v", err)
// 			msg := models.Message{
// 				Role:    "assistant",
// 				Content: "There was an error with your tool call format. Please try again with proper JSON formatting.",
// 			}
// 			messages = append(messages, msg)
// 			continue
// 		}

// 		// Unmarshal the arguments JSON string into a map
// 		var argsMap map[string]interface{}

// 		err = json.Unmarshal(jsonBytes, &argsMap)
// 		if err != nil {
// 			log.Printf("error unmarshalling arguments: %v", err)
// 			msg := models.Message{
// 				Role:    "assistant",
// 				Content: "There was an error with your tool call format. Please try again with proper JSON formatting.",
// 			}
// 			messages = append(messages, msg)
// 			continue
// 		}

// 		// Call the relevant function with parameters
// 		switch functionName {

// 		case "get_current_weather":
// 			resp := ctrl.getCurrentWeather(argsMap["location"].(string), argsMap["format"].(string))
// 			msg := models.Message{
// 				Role:    "assistant",
// 				Content: resp,
// 			}
// 			messages = append(messages, msg)
// 			continue

// 		case "get_revenue_by_month_and_year":
// 			rev, err := ctrl.db.GetRevenueByMonthYear(argsMap["month"].(int), argsMap["year"].(int))
// 			if err != nil {
// 				log.Printf("error getting revenue: %v", err)
// 				msg := models.Message{
// 					Role:    "assistant",
// 					Content: "There was an error with your tool call format. Please try again with proper JSON formatting.",
// 				}
// 				messages = append(messages, msg)
// 				continue
// 			}

// 			if rev == nil {
// 				log.Println("No revenue found for month", argsMap["month"].(int), "and year", argsMap["year"].(int))

// 				msg := models.Message{
// 					Role:    "assistant",
// 					Content: fmt.Sprintf("No revenue found for month %d and year %d", argsMap["month"].(int), argsMap["year"].(int)),
// 				}

// 				messages = append(messages, msg)
// 				continue
// 			}
// 			resp := fmt.Sprintf("Revenue for month %d and year %d is %f", rev.Month, rev.Year, rev.Amount)

// 			msg := models.Message{
// 				Role:    "assistant",
// 				Content: resp,
// 			}
// 			messages = append(messages, msg)
// 		default:
// 			log.Println("Unknown function name")
// 		}
// 	}

// 	return messages
// }

// func (ctrl *ChatController) getCurrentWeather(location string, format string) string {
// 	log.Println("Getting current weather for location", location, "in format", format)
// 	client := &http.Client{}

// 	req, err := http.NewRequest("GET", "https://wttr.in/"+location+"?format=%f", nil)
// 	log.Println("req", req)

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	resp, err := client.Do(req)

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer resp.Body.Close()

// 	var weather string

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Printf("error reading response: %v", err)
// 		return ""
// 	}
// 	weather = string(body)
// 	return fmt.Sprintf("Current weather for location %s in format %s is %s", location, format, weather)
// }

// func (ctrl *ChatController) refineToolCallResponse(toolCallResponse string, apiKey string) string {
// 	model := "microsoft/phi-3-medium-128k-instruct:free"
// 	chatBody := models.ChatBody{
// 		Model: model,
// 		Messages: []models.Message{
// 			{
// 				Role:    "user",
// 				Content: toolCallResponse,
// 			},
// 		},
// 	}

// 	responseBody, err := ctrl.callLLM(ctrl.ctx, chatBody)
// 	if err != nil {
// 		log.Printf("error calling LLM: %v", err)
// 		return ""
// 	}

// 	return responseBody.Choices[0].Message.Content
// }

func (ctrl *ChatController) GetRevenue(c *gin.Context) {
	// extract month and year from the path parameters
	month, err := strconv.Atoi(c.Param("month"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid month"})
		return
	}
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
		return
	}

	log.Printf("month: %d, year: %d", month, year)

	// get the revenue
	rev, err := ctrl.db.GetRevenueByMonthYear(month, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get revenue"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"revenue": rev.Amount})
}

func NewAgent(systemMsg string, maxRetries int, db *database.DbConfig) *models.Agent {
	return &models.Agent{
		Tools:      make(map[string]models.Tool),
		SystemMsg:  systemMsg,
		MaxRetries: maxRetries,
		Db:         db,
	}
}

func (ctrl *ChatController) GetQuarterlyRevenue(c *gin.Context) {
	// Extract quarter and year from the path parameters
	quarterStr := c.Param("quarter")
	yearStr := c.Param("year")

	quarter, err := strconv.Atoi(quarterStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quarter"})
		return
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
		return
	}

	// Get revenue for the requested quarter
	revenue, err := ctrl.GetQuarterlyRevenueInternal(quarter, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the revenue response
	c.JSON(http.StatusOK, gin.H{
		"quarter": quarter,
		"year":    year,
		"revenue": revenue,
	})
}

func (ctrl *ChatController) GetQuarterlyRevenueInternal(quarter, year int) (float64, error) {
	// Map quarters to months
	quarterMonths := map[int][]int{
		1: {1, 2, 3},
		2: {4, 5, 6},
		3: {7, 8, 9},
		4: {10, 11, 12},
	}

	months, exists := quarterMonths[quarter]
	if !exists {
		return 0, fmt.Errorf("invalid quarter: %d", quarter)
	}

	totalRevenue := 0.0
	for _, month := range months {
		revenue, err := ctrl.db.GetRevenueByMonthYear(month, year)
		if err != nil && err.Error() != fmt.Sprintf("no revenue found for month %d and year %d", month, year) {
			return 0, err
		}
		if revenue != nil {
			totalRevenue += revenue.Amount
		}
	}
	return totalRevenue, nil
}
