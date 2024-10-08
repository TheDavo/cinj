package cinj

import (
	"errors"
	"flag"
	"log"
	"os"
	"strings"

	pylex "github.com/TheDavo/cinj/lexers/python"
)

type pythonArgs struct {
	class    string
	function string
}

func newPythonArgs() *pythonArgs {
	return &pythonArgs{
		class:    "",
		function: "",
	}
}

// python uses the flag package to parse the cinj command into appropriate
// variables to later use them in the parsePython function
func (cmd CinjCommand) python() (string, error) {
	var class string
	var function string
	var content string

	pyArgs := newPythonArgs()

	pyFlag := flag.NewFlagSet("pyFlag", flag.PanicOnError)
	pyFlag.StringVar(&class, "class", "", "Grab entire content of a class")
	pyFlag.StringVar(&function, "function", "", "Grab contents of a function")

	err := pyFlag.Parse(cmd.Args)
	if err != nil {
		log.Fatal(err)
	}

	if class != "" {
		pyArgs.class = strings.Split(class, " ")[0]
	}

	if function != "" {
		pyArgs.function = strings.Split(function, " ")[0]
	}

	content, err = cmd.parsePython(*pyArgs)

	return content, err
}

// parsePython parses a python file for the appropriate content based on the
// arguments passed in the python() function call
func (cmd CinjCommand) parsePython(args pythonArgs) (string, error) {
	if args.class == "" && args.function == "" {
		content, err := cmd.returnAll()
		return content, err
	}

	content, err := os.ReadFile(cmd.Filepath)
	log.Println("Trying to use python lexer")
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	pl := pylex.NewLexer(string(content), 4)
	pl.Lex()

	// Looking only for a class
	if args.class != "" && args.function == "" {
		class, err := pl.GetClass(args.class)
		if err != nil {
			log.Fatal(err.Error())
		}

		return class, nil

	}
	// Looking for function
	if args.function != "" {
		functionText, err := pl.GetFunction(args.function, args.class)
		if err != nil {
			log.Fatal(err.Error())
		}

		return functionText, nil
	}

	return "", errors.New("Could not parse python file for wanted parameters")
}
