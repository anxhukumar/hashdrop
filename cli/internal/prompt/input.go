package prompt

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

func ReadLine(label string) (string, error) {
	fmt.Print(label)
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", errors.New("could not read input, please try again")
	}
	return strings.TrimSpace(text), nil
}

// Special function for password
func ReadPassword(label string) (string, error) {
	fmt.Print(label)
	bytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return "", errors.New("could not read password, please try again")
	}
	return string(bytes), nil
}
