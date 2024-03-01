package cinj

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Cinj struct {
	Filepath string
	Newname  string
	SrcFile  *os.File
	DestFile *os.File
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
