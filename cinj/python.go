package cinj

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type pythonArgs struct {
	decorator bool
	class     string
	function  string
}

func newPythonArgs() *pythonArgs {
	return &pythonArgs{
		decorator: false,
		class:     "",
		function:  "",
	}
}

func (cmd CinjCommand) python() (string, error) {
	var class string
	var function string
	var decorator bool
	var content string

	pyArgs := newPythonArgs()

	pyFlag := flag.NewFlagSet("pyFlag", flag.PanicOnError)
	pyFlag.StringVar(&class, "class", "", "Grab entire content of a class")
	pyFlag.StringVar(&function, "function", "", "Grab contents of a function")
	pyFlag.BoolVar(&decorator, "decorator", true,
		"Include decorators when grabbing other content")

	err := pyFlag.Parse(cmd.Args)

	if err != nil {
		log.Fatal(err)
	}

	pyArgs.decorator = decorator

	if class != "" {
		pyArgs.class = strings.Split(class, " ")[0]
	}

	if function != "" {
		pyArgs.function = strings.Split(function, " ")[0]
	}

	content, err = cmd.parsePython(*pyArgs)

	return content, err
}

func (cmd CinjCommand) parsePython(args pythonArgs) (string, error) {

	if args.class == "" && args.function == "" {
		content, err := cmd.returnAll()
		return content, err
	}

	findInsideClass := false
	foundClass := false
	lookingForClass := ""
	if args.class != "" && args.function != "" {
		findInsideClass = true
		lookingForClass = "class " + args.class
	}

	content := strings.Builder{}
	decorators := strings.Builder{}
	content.WriteString("# in file " + filepath.Base(cmd.Filepath) + "\n")

	depth := []string{}

	pythonFile, err := os.Open(cmd.Filepath)
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(pythonFile)
	if scanner.Err() != nil {
		log.Println(scanner.Err())
	}

	var lookingFor string
	if args.class != "" {
		lookingFor = "class " + args.class
	}

	if args.function != "" {
		lookingFor = "def " + args.function
	}

	index := -1
	foundToken := false
	foundEnd := false
	indexNonWs := 0 // Index of first non-whitespace character

	for scanner.Scan() {
		// var endSearch int
		// var foundLoc int
		index++
		line := scanner.Text()

		trimmedLine := strings.TrimLeft(line, " \t")
		if strings.HasPrefix(trimmedLine, "class") {
			depth = append(depth, trimmedLine)
			if findInsideClass {
				if strings.Contains(trimmedLine, lookingForClass) {
					foundClass = true
				}
			}
		}

		emptyLine := len(line) == 0
		lenTrim := 0
		if foundToken {
			lenTrim = len(line) - len(trimmedLine)
			endCriteria := lenTrim-indexNonWs <= 0
			if !emptyLine && endCriteria && !foundEnd {
				foundEnd = true
				// endSearch = index
				break
			}

			if !emptyLine {
				content.WriteString(line[indexNonWs:] + "\n")
			} else {
				content.WriteString("\n")
			}
		}

		// Start of the content has been found
		if strings.Contains(line, lookingFor) &&
			!foundToken && (findInsideClass == foundClass) {
			// foundLoc = index
			foundToken = true
			indexNonWs = len(line) - len(trimmedLine)

			// Add flavor text to provide context of the code snippet
			if len(depth) > 0 &&
				!strings.Contains(depth[len(depth)-1], lookingFor) {

				comment := depth[len(depth)-1]

				// Cuts off the colon at the end of python declarations
				if string(comment[len(comment)-1]) == ":" {
					content.WriteString("# in " +
						comment[0:len(comment)-1] + "\n")
				} else {
					content.WriteString("# in " + comment + "\n")
				}
			}

			content.WriteString(decorators.String())
			if !emptyLine {
				content.WriteString(line[indexNonWs:] + "\n")
			} else {
				content.WriteString("\n")
			}
		}

		if args.decorator {
			if strings.HasPrefix(trimmedLine, "@") {
				decorators.WriteString(trimmedLine + "\n")
			} else if !foundToken {
				decorators.Reset()
			}
		}
	}

	if !foundToken {
		log.Fatal("Token ", lookingFor, " not found!")
		return "", errors.New(fmt.Sprintf("Token %s not found!", lookingFor))
	}

	_, err = pythonFile.Seek(0, io.SeekStart)
	if err != nil {
		log.Fatal(err)
	}
	return content.String(), nil
}
