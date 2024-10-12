package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var hadError bool = false
var showTokens bool = true
var showAst bool = true
var showSource bool = false

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
	run(string(source), nil)
}

func runPrompt() {
	env := make(map[string]any)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("Error reading input: ", err)
			os.Exit(1)
		}

		source := string(line)
		source = strings.TrimSpace(source)
		if source == "exit" {
			break
		}

		if strings.HasPrefix(source, "\\set") {
			parts := strings.Split(source, " ")
			if len(parts) != 3 {
				fmt.Println("Invalid set command")
				continue
			}

			var value bool
			switch parts[2] {
			case "1":
				value = true
				break
			case "0":
				value = false
				break
			default:
				fmt.Printf("Invalid set command: \\set %s %s\n", parts[1], parts[2])
				continue
			}

			switch parts[1] {
			case "showTokens":
				showTokens = value
				break
			case "showAst":
				showAst = value
				break
			case "showSource":
				showSource = value
				break
			default:
				fmt.Printf("Invalid set command: \\set %s %s\n", parts[1], parts[2])
			}
			continue
		}

		result := run(source, &env)
		if result == nil {
			continue
		}

		if strResult, ok := result.(string); ok {
			fmt.Printf("\"%s\"\n", strResult)
		} else {
			fmt.Printf("%v\n", result)
		}
	}
}

func run(source string, env *map[string]any) any {

	if !strings.HasSuffix(source, ";") && !strings.HasSuffix(source, "}") {
		source += ";"
	}

	if showSource {
		fmt.Print(source)
	}

	//scan
	scanner := NewGloxScanner(source, reportErrorScan)
	tokens := scanner.ScanTokens()
	if hadError {
		hadError = false
		return nil
	}
	if showTokens {
		fmt.Println(tokens)
	}

	//parse
	parser := NewParser(tokens, reportErrorParse)
	stmts, err := parser.parse()
	if err != nil {
		fmt.Println("Error parsing expression: ", err)
		return nil
	}

	if showAst {
		astPrinter := NewAstPrinter(env)
		err = astPrinter.print(stmts)
		if err != nil {
			fmt.Println("Error printing tree: ", err)
			return nil
		}
	}

	// run
	interp := NewInterpreter(env)
	val, err := interp.interpert(stmts)
	if err != nil {
		fmt.Println("Error interpreting: ", err)
		return nil
	}
	return val
}
