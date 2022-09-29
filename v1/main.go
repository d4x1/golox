package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var (
	disableDebugScanner       = true
	disableDebugResolveLocals = true
)

var hasErr bool

func run(source string) error {
	scanner := newScanner(source)
	tokens, err := scanner.scanTokens()
	if err != nil {
		return err
	}
	if !disableDebugScanner {
		fmt.Println(strings.ToUpper("[debug scanner]"))
		for _, token := range tokens {
			fmt.Println(token)
		}
	}

	parser := newParser(tokens)
	stmts, err := parser.parse()
	if err != nil {
		return err
	}

	fmt.Println(strings.ToUpper("[debug execute stmts]"))
	intp := newInterpreter()
	resolver := newResolver(intp)
	if err := resolver.resolveStmts(stmts); err != nil {
		return err
	}
	if !disableDebugResolveLocals {
		fmt.Printf("interpreter locals: %+v\n", intp.locals)
	}
	intp.interpret(stmts)

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
	if err := run(string(bytes)); err != nil {
		return err
	}
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
	// fmt.Println("args:", lenArgs)
	if lenArgs > 2 {
		os.Exit(1)
	} else if lenArgs == 2 {
		if err := runFile(args[1]); err != nil {
			fmt.Println("Run File Error: ", err)
		}
	} else {
		if err := runPrompt(); err != nil {
			fmt.Println("Error: ", err)
		}
	}
}
