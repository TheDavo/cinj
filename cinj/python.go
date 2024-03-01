package cinj

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func (cmd CinjCommand) parsePython(token string, tokenValue string) string {

	content := strings.Builder{}
	content.WriteString("# in file " + filepath.Base(cmd.Filepath) + "\n")

	depth := []string{}

	pythonFile, err := os.Open(cmd.Filepath)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(pythonFile)
	if scanner.Err() != nil {
		log.Println(scanner.Err())
	}

	var lookingFor string
	switch token {
	case "class":
		lookingFor = "class " + tokenValue
	case "function":
		lookingFor = "def " + tokenValue
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
		}

		lenTrim := 0
		if foundToken {
			emptyLine := len(line) == 0
			lenTrim = len(line) - len(strings.TrimLeft(line, " \t"))
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

		if strings.Contains(line, lookingFor) && !foundToken {
			// foundLoc = index
			foundToken = true
			indexNonWs = len(line) - len(strings.TrimLeft(line, " \t"))

			// Add flavor text to provide context of the code snippet
			if len(depth) > 0 && !strings.Contains(depth[len(depth)-1], lookingFor) {
				comment := depth[len(depth)-1]

				// Cuts off the colon at the end of python declarations
				if string(comment[len(comment)-1]) == ":" {
					content.WriteString("# in " + comment[0:len(comment)-1] + "\n")
				} else {
					content.WriteString("# in " + comment + "\n")
				}
			}
			if len(line) != 0 {
				content.WriteString(line[indexNonWs:] + "\n")
			} else {
				content.WriteString("\n")
			}
		}
	}

	if !foundToken {
		log.Fatal("Token ", lookingFor, " not found!")
	}

	_, err = pythonFile.Seek(0, io.SeekStart)
	if err != nil {
		log.Fatal(err)
	}
	return content.String()
}

func (cmd CinjCommand) python() []string {
	var class string
	var function string
	var decorator bool
	var content []string

	pyFlag := flag.NewFlagSet("pyFlag", flag.PanicOnError)
	pyFlag.StringVar(&class, "class", "", "Grab entire content of a class")
	pyFlag.StringVar(&function, "function", "", "Grab contents of a function")
	pyFlag.BoolVar(&decorator, "decorator", true, "Include decorators when grabbing other content")

	err := pyFlag.Parse(cmd.Args)

	if err != nil {
		log.Fatal(err)
	}

	if class != "" {
		classes := strings.Split(class, " ")
		for _, c := range classes {
			content = append(content, cmd.parsePython("class", c))
		}
	}

	if function != "" {
		functions := strings.Split(function, " ")
		for _, f := range functions {
			content = append(content, cmd.parsePython("function", f))
		}
	}

	return content
}
