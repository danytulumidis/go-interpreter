package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/lexer"
	"monkey/parser"
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
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}

const MONKEY_FACE = `            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
