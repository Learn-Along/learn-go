package codegenerator

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/term"
)

type Eliql struct {
	hadError bool
}

// RunFile Runs the code found in the given file
func (q *Eliql) RunFile(filePath string) error {
	source, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return q.run(string(source))
}

// Repl Runs the EliQL in a REPL
func (q *Eliql) Repl() error {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("initializing terminal error: %s", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	terminal := term.NewTerminal(os.Stdin, ">")

	for {
		source, err := terminal.ReadLine()
		if err != nil {
			return err
		}

		q.run(source)
	}
}

func (q *Eliql) run(source string) error {
	scanner := NewScanner(q, source)
	tokens, err := scanner.ScanTokens()
	if err != nil {
		return err
	}

	for _, token := range tokens {
		fmt.Printf("%v\n", token)
	}

	return nil
}

func (q *Eliql) Error(line int64, message string) {
	q.report(line, "", message)
}

func (q *Eliql) report(line int64, where string, message string) {
	log.Printf("[line %d] Error %s: %s", line, where, message)
	q.hadError = true
}
