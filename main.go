package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	gpt3 "github.com/PullRequestInc/go-gpt3"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func GetResponse(client gpt3.Client, ctx context.Context, quesiton string) {
	err := client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: []string{
			quesiton,
		},
		MaxTokens:   gpt3.IntPtr(3000),
		Temperature: gpt3.Float32Ptr(0),
	}, func(resp *gpt3.CompletionResponse) {
		fmt.Print(resp.Choices[0].Text)
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(13)
	}
	fmt.Printf("\n")
}

type NullWriter int

func (NullWriter) Write([]byte) (int, error) { return 0, nil }

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.SetOutput(new(NullWriter))
	apiKey := os.Getenv("CHATGP_API_KEY")
	if apiKey == "" {
		panic("CHATGP_API_KEY isn't set")
	}

	ctx := context.Background()
	client := gpt3.NewClient(apiKey)
	rootCmd := &cobra.Command{
		Use:   "chatgpt",
		Short: "Chat with ChatGPT in console.",
		Run: func(cmd *cobra.Command, args []string) {
			scanner := bufio.NewScanner(os.Stdin)
			quit := false

			for !quit {
				fmt.Print(">> ")

				if !scanner.Scan() {
					break
				}

				question := scanner.Text()
				questionParam := validateQuestion(question)
				switch questionParam {
				case "quit":
					quit = true
				case "":
					continue

				default:
					GetResponse(client, ctx, questionParam)
				}
			}
		},
	}

	log.Fatal(rootCmd.Execute())
}

func validateQuestion(question string) string {
	quest := strings.Trim(question, " ")
	keywords := []string{"", "loop", "break", "cls", "exit", "block"}
	for _, x := range keywords {
		if quest == x {
			return ""
		}
	}
	return quest
}
