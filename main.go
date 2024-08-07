package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

var hadError bool = false

func main() {
	if len(os.Args) > 2 {
		fmt.Println("usage: glox [file]")
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}

func runFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening file %s: %s\n", path, err)
		os.Exit(1)
	}
	defer f.Close()
	source, err := io.ReadAll(f)
	run(string(source))
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("Error reading input: ", err)
			os.Exit(1)
		}
		run(string(line))
	}
}

func run(source string) {
	fmt.Print(source)
	scanner := NewGloxScanner(source, reportErrorScan)
	tokens := scanner.ScanTokens()
	fmt.Println(tokens)
	if hadError {
		hadError = false
		return
	}
	parser := NewParser(tokens, reportErrorParse)
	expr, err := parser.parse()
	if err != nil {
		fmt.Println("Error parsing expression: ", err)
		return
	}
	astPrinter := AstPrinter{}
	treeStr, err := expr.accept(astPrinter)
	if err != nil {
		fmt.Println("Error printing tree: ", err)
		return
	}
	fmt.Println(treeStr)

	interp := Interpreter{}
	interp.interpert(expr)
}
