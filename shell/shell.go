package shell

import (
	"log"
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
	log.Printf("shell - initing print")
	for {
		select {
		case evt := <- c.chatClient.OutChan:
			if evt.ClientEventType == client.ClientEventType_StopProcessing {
				log.Printf("shell - quitting")
				os.Exit(0)
			}

			text, ok := evt.Payload.(string)
			if !ok {
				continue
			}		
			fmt.Printf("< %s\n", text)
		}
	}
}

func (c *Shell) readLines() {
	log.Printf("shell - initing readLines")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()

		if text == "quit" {
			log.Printf("shell - quitting")
			os.Exit(0)
		}

		fmt.Printf("> %s\n", text)
		c.chatClient.InChan <- &client.ClientEvent {
			ClientEventType: client.ClientEventType_Message,
			Payload: text,
		}
	}
}