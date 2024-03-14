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

func (c *Cinj) Run() error {
	file, err := os.Open(c.Filepath)
	if err != nil {
		fmt.Println("Failure opening file", c.Filepath)
		log.Fatal(err)
	}
	defer file.Close()

	newFile, err := os.Create(c.Newname)
	if err != nil {
		fmt.Print(err)
		return err
	}
	defer file.Close()

	c.SrcFile = file
	c.DestFile = newFile

	err = c.cinj()
	if err != nil {
		newFile.Close()
		os.Remove(c.Newname)
		return err
	}

	newFile.Close()
	return nil
}

func (c *Cinj) cinj() error {

	srcScanner := bufio.NewScanner(c.SrcFile)

	for srcScanner.Scan() {
		line := srcScanner.Text()

		if strings.HasPrefix(line, "cinj") {
			command, err := c.getCinjCommand(line)
			if err != nil {
				log.Fatal(err)
				return err
			}

			language := command.fileExtForMarkDown()
			content, err := c.getContentFromCommand(command)
			if err != nil {
				return err
			}

			contentScanner := bufio.NewScanner(strings.NewReader(content))

			c.DestFile.WriteString("```" + language.String() + "\n")
			for contentScanner.Scan() {
				contentLine := contentScanner.Text()
				_, err := c.DestFile.WriteString(contentLine + "\n")

				if err != nil {
					fmt.Println(err)
					return err
				}
			}
			c.DestFile.WriteString("```\n")
			srcScanner.Scan()

		} else {
			c.DestFile.WriteString(line + "\n")
		}
	}
	return nil
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

func (c Cinj) getContentFromCommand(cmd CinjCommand) (string, error) {
	switch cmd.FileType {
	case Python:
		content, err := cmd.python()
		return content, err
	default:
		content, err := cmd.returnAll()
		return content, err
	}
}
