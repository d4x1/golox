package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var hasErr bool

func run(source string) error {
	scanner := newScanner(source)
	tokens, err := scanner.scanTokens()
	if err != nil {
		return err
	}
	fmt.Println(strings.ToUpper("[debug scanner]"))
	for _, token := range tokens {
		fmt.Println(token)
	}

	parser := newParser(tokens)
	expression, err := parser.parse()
	if err != nil {
		return err
	}

	fmt.Println(strings.ToUpper("[debug parser expression]"))
	expressionString := expression.acceptStringVisitor(&PrettyPrinter{})
	fmt.Println(expressionString)

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
		fmt.Print("Input: ", line)
		if err := run(line); err != nil {
			hasErr = true
			fmt.Printf("Error: %+v\n", err)
		} else {
			hasErr = false
		}
	}
}
func main() {
	fmt.Println(strings.ToUpper("welcome to go lox!"))
	args := os.Args
	lenArgs := len(args)
	if lenArgs > 2 {
		os.Exit(1)
	} else if lenArgs == 2 {
		runFile(args[1])
	} else {
		if err := runPrompt(); err != nil {
			fmt.Println("Error: ", err)
		}
	}
}
