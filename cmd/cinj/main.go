package main

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

type Filetype string

const (
	Python     Filetype = "python"
	Javascript          = "javascript"
	Markdown            = "md"
	Text                = ""
	Plain               = ""
)

type SupportedArgs map[Filetype][]string

func (ft Filetype) String() string {
	return string(ft)
}

type PythonArgs struct {
	Class      []string
	Func       []string
	Decorators bool
}

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

func main() {
	var cinj Cinj
	var fp string
	var newname string

	flag.StringVar(&fp, "fp", "", "Filepath of the Markdown file to Cinj")
	flag.StringVar(&newname, "newname", "", "New name for output file")

	flag.Parse()
	if fp == "" {
		log.Fatal("No file provided")
		os.Exit(1)
	}

	absFp, err := filepath.Abs(fp)

	if err != nil {
		log.Fatal(err)
	}

	fileext := filepath.Ext(absFp)

	if fileext != ".md" && fileext != ".cinj" {
		log.Fatal("Unrecognized or incorrect filetype")
		os.Exit(1)
	}
	cinj.Filepath = absFp

	// Default case
	if newname == "" {
		temp, found := strings.CutSuffix(absFp, fileext)
		if found && fileext == ".md" {
			cinj.Newname = temp + ".cinj.md"
		} else {
			cinj.Newname = temp + ".md"
		}
	} else {
		if filepath.IsAbs(newname) {
			// Handle the case that the newname flag is the same as the original file
			if newname == absFp {
				if fileext == ".md" {
					newname = newname[:len(newname)-len(fileext)] + ".cinj.md"
				} else {
					newname = newname[:len(newname)-len(fileext)] + ".md"
				}
				cinj.Newname = newname
			}
			cinj.Newname = newname
		} else {
			absNew, _ := filepath.Abs(newname)
			if absNew == absFp {
				if fileext == ".md" {
					absNew = absNew[:len(absNew)-len(fileext)] + ".cinj.md"
				} else {
					absNew = absNew[:len(absNew)-len(fileext)] + ".md"
				}
			}
			cinj.Newname = absNew
		}
	}

	fmt.Println("newname", cinj.Newname)
	cinj.Run()
}

func (c *Cinj) Run() {
	file, err := os.Open(c.Filepath)
	if err != nil {
		fmt.Println("Failure opening file", c.Filepath)
		log.Fatal(err)
	}
	defer file.Close()

	newFile, err := os.Create(c.Newname)
	if err != nil {
		log.Fatal(err)
	}
	defer newFile.Close()

	c.SrcFile = file
	c.DestFile = newFile

	c.cinj()
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
		fmt.Println(cmd.Args)
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

	fmt.Println(lookingFor)

	index := -1
	foundToken := false
	foundEnd := false
	indexNonWs := 0 // Index of first non-whitespace character

	for scanner.Scan() {
		// var endSearch int
		// var foundLoc int
		index++
		line := scanner.Text()

		lenTrim := 0
		if foundToken {
			emptyLine := len(line) == 0
			lenTrim = len(line) - len(strings.TrimLeft(line, " \t"))
			endCriteria := lenTrim-indexNonWs <= 0
			content.WriteString(line + "\n")
			if !emptyLine && endCriteria && !foundEnd {
				foundEnd = true
				// endSearch = index
			}
		}

		if strings.Contains(line, lookingFor) && !foundToken {
			// foundLoc = index
			foundToken = true
			indexNonWs = len(line) - len(strings.TrimLeft(line, " \t"))
			content.WriteString(line + "\n")
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
