package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/lexer"
	"monkey/token"
)

const PROMPT = ">> "

// io Reader = User Input from the Console
// io Writer = Output to the console
func Start(in io.Reader, out io.Writer) {
	// Scanner for reading User Input
	scanner := bufio.NewScanner(in)

	// Endless loop
	for {
		// Writes PROMPT into the out (io.Writer) so to the console
		fmt.Fprintf(out, PROMPT)
		// Scans the user Input, when nothing is there (user pressed cmd+c) exit loop and REPL
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		// Get the user Input as a string
		line := scanner.Text()
		// Init the Lexer
		l := lexer.New(line)

		// Reads all the tokens one by one until the end of the file and prints them
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}
