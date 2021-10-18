package server

import (
	"fmt"
	"strings"
	"bufio"
)

type nameCommand struct {
	name string
}

type messageCommand struct {
	target string
	message string
}

// name [name]
// message [all | name] [message]
func parseToCommand(input string) (interface{}, error) {
	reader := bufio.NewReader(strings.NewReader(fmt.Sprintf("%s\n", input)))
	cmd, err := reader.ReadString(' ')
	if err != nil {
		return nil, err
	}
	cmd = cmd[0:len(cmd)-1]

	switch cmd {
	case "name":
		name, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		name = name[:len(name)-1]	
		return &nameCommand {
			name: name,
		}, nil
	case "message":
		target, err := reader.ReadString(' ')
		if err != nil {
			return nil, err
		}

		message, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		message = message[:len(message)-1]

		return &messageCommand {
			target: target,
			message: message,
		}, nil
	default:
		return nil, fmt.Errorf("Command not found")
	}	
}