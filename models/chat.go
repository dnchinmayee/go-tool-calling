package models

import "tempfunctiontools/internal/database"

// Chat message role defined by the OpenAI API.
const (
	ChatMessageRoleSystem    = "system"
	ChatMessageRoleUser      = "user"
	ChatMessageRoleAssistant = "assistant"
	ChatMessageRoleFunction  = "function"
	ChatMessageRoleTool      = "tool"
	ChatMessageRoleDeveloper = "developer"
)

type ChatBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Tools    []Tool    `json:"tools,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role      string     `json:"role"`
			Content   string     `json:"content"`
			ToolCalls []ToolCall `json:"tool_calls"` //list of tool calls
		} `json:"message"`
	} `json:"choices"`

	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// after u have added your tools to chatBody
type ToolCall struct {
	Index    int    `json:"index"`
	Id       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	}
}

type UserResponse struct {
	Messages    []Message    `json:"messages"`     //returns chatBody messages(role and content)
	LLMResponse ChatResponse `json:"llm_response"` //add LLM response
}

//"tools": [
// {
// "type": "function",
//   "function": {
//     "name": "

type Tool struct {
	Function *Function                              `json:"function"`
	Type     string                                 `json:"type"`
	Execute  func(args map[string]any) (any, error) `json:"-"`
}

type Function struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Parameters  *Parameters `json:"parameters"`
}

type Parameters struct {
	Type       string                `json:"type"`
	Properties map[string]*Parameter `json:"properties"`
	Required   []string              `json:"required"`
}

type Parameter struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
}

type Agent struct {
	Tools      map[string]Tool
	SystemMsg  string
	MaxRetries int
	Db         *database.DbConfig
}
