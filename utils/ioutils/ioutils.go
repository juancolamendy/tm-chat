package ioutils

import (
	"fmt"
	"bufio"
)

func SendString(out *bufio.Writer, text string) error {
	_, err := out.WriteString(fmt.Sprintf("%s\n",text))
	if err != nil {
		return err
	}
	return out.Flush()
}

func ReceiveString(in *bufio.Reader) (string, error) {
	text, err := in.ReadString('\n')
	if err != nil {
		return "", err
	}
	return text[:len(text)-1], nil
}