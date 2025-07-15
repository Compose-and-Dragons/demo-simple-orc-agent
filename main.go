package main

import (
	"fmt"
	"os"

	"github.com/budgies-nest/budgie/agents"
	"github.com/budgies-nest/budgie/helpers"
	"github.com/openai/openai-go"
)

func main() {

	modelRunnerBaseUrl := os.Getenv("MODEL_RUNNER_BASE_URL")

	if modelRunnerBaseUrl == "" {
		panic("MODEL_RUNNER_BASE_URL environment variable is not set")
	}
	modelRunnerChatModel := os.Getenv("MODEL_RUNNER_CHAT_MODEL")
	fmt.Println("Using Model Runner Chat Model:", modelRunnerChatModel)

	if modelRunnerChatModel == "" {
		panic("MODEL_RUNNER_CHAT_MODEL environment variable is not set")
	}

	systemInstruction, err := helpers.ReadTextFile("instructions.md")
	if err != nil {
		panic(err)
	}
	characterSheet, err := helpers.ReadTextFile("character_sheet.md")
	if err != nil {
		panic(err)
	}

	// Create a new agent named Bob
	npcAgent, err := agents.NewAgent("Bob",
		agents.WithDMR(modelRunnerBaseUrl),
		agents.WithParams(openai.ChatCompletionNewParams{
			Model:       modelRunnerChatModel,
			Temperature: openai.Opt(0.8),
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage("CONTEXT:\n" + characterSheet),
				openai.SystemMessage(systemInstruction),
			},
		}),
	)
	if err != nil {
		panic(err)
	}

	// Start the TUI prompt with custom messages
	err = npcAgent.Prompt(agents.PromptConfig{
		UseStreamCompletion:        true, // Set to false for non-streaming completion
		StartingMessage:            "👹 I'm an Orc",
		ExplanationMessage:         "Ask me anything about me. Type '/bye' to quit or Ctrl+C to interrupt responses.",
		PromptTitle:                "✋ Query",
		ThinkingPrompt:             "⏳",
		InterruptInstructions:      "(Press Ctrl+C to interrupt)",
		CompletionInterruptMessage: "⚠️ Response was interrupted\n",
		GoodbyeMessage:             "👹 Bye!",
	})
	if err != nil {
		panic(err)
	}
}
