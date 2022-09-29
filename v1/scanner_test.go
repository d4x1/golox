package main

import (
	"fmt"
	"testing"
)

func Test_scanner_scanTokens(t *testing.T) {
	raw := "!*+-/=<> <= == \"golox\""
	tokens, err := newScanner(raw).scanTokens()
	if err != nil {
		t.Error(err)
	} else {
		for _, token := range tokens {
			fmt.Println(token)
		}
	}
}
