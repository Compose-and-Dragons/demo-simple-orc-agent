package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"orc-agent/helpers"
	"github.com/charmbracelet/huh"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
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

	systemInstructions, err := helpers.ReadTextFile("instructions.md")
	if err != nil {
		panic(err)
	}
	characterSheet, err := helpers.ReadTextFile("character_sheet.md")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	clientEngine := openai.NewClient(
		option.WithBaseURL(modelRunnerBaseUrl),
		option.WithAPIKey(""),
	)

	// Chat Completion parameters
	chatCompletionParams := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("CONTEXT:\n" + characterSheet),
			openai.SystemMessage(systemInstructions),
		},
		Model:       modelRunnerChatModel,
		Temperature: openai.Opt(0.8),
	}

	type PromptConfig struct {
		StartingMessage            string
		ExplanationMessage         string
		PromptTitle                string
		ThinkingPrompt             string
		InterruptInstructions      string
		CompletionInterruptMessage string
		GoodbyeMessage             string
	}
	promptConfig := PromptConfig{
		StartingMessage:       "👹 I'm an Orc",
		ExplanationMessage:    "Ask me anything about me. Type '/bye' to quit or Ctrl+C to interrupt responses.",
		PromptTitle:           "✋ Query",
		ThinkingPrompt:        "⏳",
		InterruptInstructions: "(Press Ctrl+C to interrupt)",
		//CompletionInterruptMessage: "⚠️ Response was interrupted\n",
		GoodbyeMessage: "👹 Bye!",
	}

	fmt.Println(promptConfig.StartingMessage)
	fmt.Println(promptConfig.ExplanationMessage)

	for {
		fmt.Print(promptConfig.ThinkingPrompt)
		fmt.Println(promptConfig.InterruptInstructions)

		var userInput string

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewText().
					Title(promptConfig.PromptTitle).
					Placeholder("Type your question here...").
					Value(&userInput).
					ExternalEditor(false),
			),
		)

		// Run the form
		err := form.Run()
		if err != nil {
			// TODO: handle error
		}

		// Trim whitespace
		userInput = strings.TrimSpace(userInput)

		// Check for empty input
		if userInput == "" {
			continue
		}

		// Check for /bye command
		if userInput == "/bye" {
			fmt.Println(promptConfig.GoodbyeMessage)
			break
		}


		fmt.Println("🤖 Starting chat completion...")
		fmt.Println(strings.Repeat("=", 80))

		// CHAT COMPLETION:
		chatCompletionParams.Messages = append(
			chatCompletionParams.Messages,
			openai.UserMessage(userInput),
		)

		stream := clientEngine.Chat.Completions.NewStreaming(ctx, chatCompletionParams)

		for stream.Next() {
			chunk := stream.Current()
			// Stream each chunk as it arrives
			if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
				fmt.Print(chunk.Choices[0].Delta.Content)
			}
		}

		if err := stream.Err(); err != nil {
			fmt.Printf("😡 Stream error: %v\n", err)
		}
		
		fmt.Println()
		fmt.Println(strings.Repeat("=", 80))
		fmt.Println() // Add spacing between interactions
	}

}
