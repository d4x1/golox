package main

import (
	"bufio"
	"fmt"
	"os"
)

var hasErr bool

func run(source string) error {
	scanner := newScanner(source)
	tokens, err := scanner.scanTokens()
	if err != nil {
		return err
	}
	for _, token := range tokens {
		fmt.Println(token)
	}
	return nil
}

func printError(line int, msg string) {
	fmt.Printf("line: %d, err: %s\n", line, msg)
	hasErr = true
}

func runFile(fileName string) error {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	run(string(bytes))
	if hasErr {
		os.Exit(1)
	}
	return nil
}

func runPrompt() error {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("golox > ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		fmt.Print(line)
		run(line)
		hasErr = false
	}
}
func main() {
	fmt.Println("welcome to go lox")
	args := os.Args
	lenArgs := len(args)
	if lenArgs > 2 {
		os.Exit(1)
	} else if lenArgs == 2 {
		runFile(args[1])
	} else {
		runPrompt()
	}

}
