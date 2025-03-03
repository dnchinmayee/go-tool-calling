package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"tempfunctiontools/internal/functions"
	"tempfunctiontools/models"
)

func (ctrl *ChatController) ProcessQuery(ctx context.Context, chatBody models.ChatBody) ([]models.Message, error) {
	// add tool calls
	chatBody.Tools = functions.GetTools(ctrl.agent)

	messages := chatBody.Messages

	initialResponse, err := ctrl.createInitialCompletion(ctx, chatBody)
	if err != nil {
		return nil, err
	}

	log.Printf("initialResponse: %+v", initialResponse)

	// if choices is empty, return error
	if len(initialResponse.Choices) == 0 {
		log.Printf("no choices in initial response")
		msg := models.Message{
			Role:    models.ChatMessageRoleAssistant,
			Content: "Error: No response from LLM",
		}

		messages = append(messages, msg)
		return messages, nil
		// return messages, fmt.Errorf("no choices in initial response")
	}

	initialMsg := models.Message{
		Role:    initialResponse.Choices[0].Message.Role,
		Content: initialResponse.Choices[0].Message.Content,
	}
	// if no tool calls, add initial response
	if len(initialResponse.Choices[0].Message.ToolCalls) == 0 {
		messages = append(messages, initialMsg)
		return messages, nil
	}

	// execute tool calls
	var toolResults []models.Message

	toolCalls := initialResponse.Choices[0].Message.ToolCalls
	for _, toolCall := range toolCalls {
		result, err := ctrl.executeToolCall(ctx, toolCall)
		if err != nil {
			// return nil, err
			result = models.Message{
				Role:    models.ChatMessageRoleUser,
				Content: "Error: " + err.Error(),
			}
		}
		toolResults = append(toolResults, result)
	}

	// create final response
	finalResponse, err := ctrl.createFinalResponse(ctx, chatBody, initialMsg, toolResults)
	if err != nil {
		return nil, err
	}

	if len(finalResponse.Choices) == 0 {
		log.Printf("no choices in final response")
		msg := models.Message{
			Role:    models.ChatMessageRoleAssistant,
			Content: "No response from LLM",
		}
		messages = append(messages, msg)
		return messages, fmt.Errorf("no choices in final response")
	}

	// add final response
	finalResponseMsg := models.Message{
		Role:    finalResponse.Choices[0].Message.Role,
		Content: finalResponse.Choices[0].Message.Content,
	}
	messages = append(messages, finalResponseMsg)

	return messages, nil
}

// createInitialCompletion sends the query to the LLM and gets the initial response
func (ctrl *ChatController) createInitialCompletion(ctx context.Context, chatBody models.ChatBody) (models.ChatResponse, error) {
	resp, err := ctrl.callLLM(ctx, chatBody)
	if err != nil {
		log.Printf("error calling LLM: %v", err)
		return models.ChatResponse{}, err
	}

	return resp, nil
}

// executeToolCall sends the tool call to the LLM and gets the response
func (ctrl *ChatController) executeToolCall(ctx context.Context, toolCall models.ToolCall) (models.Message, error) {
	// get function name and args
	functionName := toolCall.Function.Name
	arguments := toolCall.Function.Arguments

	log.Printf("Function name: %s, Arguments: %v", functionName, arguments)
	var args map[string]any
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		log.Printf("error unmarshalling arguments: %v", err)
		return models.Message{}, err
	}

	log.Printf("Arguments: %+v", args)
	// call function
	tool, exists := ctrl.agent.Tools[functionName]
	if !exists {
		log.Printf("tool %s not found", functionName)
		return models.Message{}, fmt.Errorf("tool %s not found", functionName)
	}

	result, err := tool.Execute(args)
	if err != nil {
		log.Printf("error executing tool: %v", err)
		return models.Message{}, err
	}

	// convert to json
	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Printf("error marshalling response: %v", err)
		return models.Message{}, err
	}

	resp := models.Message{
		Role:    models.ChatMessageRoleUser,
		Content: string(resultJSON),
	}

	return resp, nil
}

// create final response, based on chatbody, initial response and tool results
func (ctrl *ChatController) createFinalResponse(ctx context.Context, chatBody models.ChatBody, initialMessage models.Message, toolResults []models.Message) (models.ChatResponse, error) {
	// check chat body if it has tools then remove them
	if len(chatBody.Tools) > 0 {
		chatBody.Tools = nil
	}

	// add initial response and tool results
	if initialMessage.Content != "" {
		chatBody.Messages = append(chatBody.Messages, initialMessage)
	}
	chatBody.Messages = append(chatBody.Messages, toolResults...)

	// call LLM
	response, err := ctrl.callLLM(ctx, chatBody)
	if err != nil {
		log.Printf("error calling LLM: %v", err)
		return models.ChatResponse{}, err
	}

	return response, nil
}
