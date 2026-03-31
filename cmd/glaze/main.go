package main

import (
	"fmt"
	"os"
	"glaze/internal/evaluator"
	"glaze/internal/lexer"
	"glaze/internal/object"
	"glaze/internal/parser"
	"glaze/internal/token"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No input file provided.")
		fmt.Println("Usage: ./glaze <filename>")
		os.Exit(1)
	}
	if len(os.Args) > 2 {
		fmt.Println("Multiple files found.")
		fmt.Println("Usage: ./glaze <filename>")
		os.Exit(1)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	source := string(data)

	// ==========================================
	// 1. LEXING (Text -> Tokens)
	// ==========================================
	s := lexer.InitScanner(source)

	var tokens []token.Token
	for {
		// NOTE: Change `NextToken()` to whatever your lexer uses to get the next token
		tok := s.ScanToken()
		tokens = append(tokens, tok)

		// NOTE: Change `token.EOF` if your end-of-file constant is named differently
		if tok.Type == token.END_OF_FILE {
			break
		}
	}

	// ==========================================
	// 2. PARSING (Tokens -> AST)
	// ==========================================
	p := parser.New(tokens)
	program := p.ParseProgram()

	// Stop execution if there are syntax errors!
	if len(p.Errors()) != 0 {
		fmt.Println("PARSER ERRORS:")
		for _, msg := range p.Errors() {
			fmt.Printf("\t- %s\n", msg)
		}
		os.Exit(1)
	}

	// ==========================================
	// 3. EVALUATING (AST -> Execution)
	// ==========================================
	env := object.NewEnvironment()
	result := evaluator.Eval(program, env)

	// Catch and print any runtime type-mismatches or undeclared variable errors
	if result != nil && result.Type() == object.ERROR_OBJ {
		fmt.Println("RUNTIME ERROR:")
		fmt.Println(result.Inspect())
		os.Exit(1)
	}
}
