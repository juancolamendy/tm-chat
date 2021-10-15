package shell

import (
	"bufio"
	"os"
	"fmt"

	"github.com/juancolamendy/tm-chat/client"
)

type Shell struct {
	chatClient *client.ChatClient
}

func NewShell(chatClient *client.ChatClient) *Shell {
	return &Shell {
		chatClient: chatClient,
	}
}

func (c *Shell) Init() {
	go c.print()

	c.readLines()
}

func (c *Shell) print() {
	for {
		select {
		case text, ok := <- c.chatClient.OutChan:
			if !ok {
				return
			}
			fmt.Printf("< %s\n", text)
		}
	}
}

func (c *Shell) readLines() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Printf("> %s\n", text)
		c.chatClient.InChan <- text
	}
}