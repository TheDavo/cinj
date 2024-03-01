package cinj

import (
	"bufio"
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Filetype string

func (ft Filetype) String() string {
	return string(ft)
}

const (
	Python     Filetype = "python"
	Javascript          = "javascript"
	Markdown            = "md"
	Text                = ""
	Plain               = ""
)

type Cinj struct {
	Filepath string
	Newname  string
	SrcFile  *os.File
	DestFile *os.File
}

type CinjCommand struct {
	Filepath string
	Args     []string
	FileType Filetype
	SuppArgs []string
}

func (c *Cinj) cinj() {

	srcScanner := bufio.NewScanner(c.SrcFile)

	for srcScanner.Scan() {
		line := srcScanner.Text()

		if strings.HasPrefix(line, "cinj") {
			command, err := c.getCinjCommand(line)
			if err != nil {
				log.Fatal(err)
			}

			language := command.fileExtForMarkDown()
			contents := c.getContentFromCommand(command)

			for _, content := range contents {
				contentScanner := bufio.NewScanner(strings.NewReader(content))

				c.DestFile.WriteString("```" + language.String() + "\n")
				for contentScanner.Scan() {
					contentLine := contentScanner.Text()
					_, err := c.DestFile.WriteString(contentLine + "\n")

					if err != nil {
						log.Fatal(err)
					}
				}
				c.DestFile.WriteString("```\n")
				srcScanner.Scan()

			}
		} else {
			c.DestFile.WriteString(line + "\n")
		}
	}
}

func (c Cinj) getCinjCommand(s string) (CinjCommand, error) {
	var cmd CinjCommand
	if len(s) <= 6 {
		return cmd, errors.New("Cinj command found too short, must contain 'cinj{arg}' at minimum")
	}

	content := s[5 : len(s)-1]
	contentSplit := strings.Split(content, " ")

	cmd.Filepath = contentSplit[0]
	if filepath.IsLocal(cmd.Filepath) {
		resolvedPath := filepath.Join(filepath.Dir(c.Filepath), cmd.Filepath)
		cmd.Filepath = resolvedPath
		cmd.FileType = cmd.fileExtForMarkDown()
	}
	if len(contentSplit) > 1 {
		cmd.Args = contentSplit[1:]
	}

	return cmd, nil
}

func (c Cinj) getContentFromCommand(cmd CinjCommand) []string {
	switch cmd.FileType {
	case Python:
		return cmd.python()
	default:
		return []string{cmd.returnAll()}
	}
}

func (cmd CinjCommand) extractContent() ([]string, error) {
	return []string{}, nil
}

func (cmd CinjCommand) returnAll() string {
	content, err := os.ReadFile(cmd.Filepath)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	return string(content)

}

func (cmd CinjCommand) fileExtForMarkDown() Filetype {
	switch filepath.Ext(cmd.Filepath) {
	case ".py":
		return Python
	case ".js":
		return Javascript
	case ".txt":
		return Text
	case ".md":
		return Markdown
	}

	return Plain
}

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
