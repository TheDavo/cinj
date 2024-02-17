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

type Filetype int

const (
	Python Filetype = iota
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

			syntaxHighlight := c.getFileExtForMarkdown(command)
			content := c.getContentFromCommand(command)
			contentScanner := bufio.NewScanner(strings.NewReader(content))

			c.DestFile.WriteString("```" + syntaxHighlight + "\n")
			for contentScanner.Scan() {
				contentLine := contentScanner.Text()
				_, err := c.DestFile.WriteString(contentLine + "\n")

				if err != nil {
					log.Fatal(err)
				}
			}
			c.DestFile.WriteString("```\n")
			srcScanner.Scan()
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
	}
	if len(contentSplit) > 1 {
		cmd.Args = contentSplit[1:]
		fmt.Println(cmd.Args)
	}

	return cmd, nil
}

func (c Cinj) getContentFromCommand(cmd CinjCommand) string {
	content, err := os.ReadFile(cmd.Filepath)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	return string(content)
}

func (c Cinj) getFileExtForMarkdown(cmd CinjCommand) string {
	switch filepath.Ext(cmd.Filepath) {
	case ".py":
		return "python"
	case ".js":
		return "javascript"
	}

	return ""
}

func ParsePy(file *os.File, cmd CinjCommand) {

	scanner := bufio.NewScanner(file)

	if scanner.Err() != nil {
		log.Println(scanner.Err())
	}

	token := "Test2"
	index := -1
	found_token := false
	found_end := false
	found_loc := 0
	index_non_ws := 0 // Index of first non-whitespace character
	end_search := 0

	for scanner.Scan() {
		index++
		line := scanner.Text()

		len_trim := 0
		if found_token {
			empty_line := len(line) == 0
			len_trim = len(line) - len(strings.TrimLeft(line, " \t"))
			end_critera := len_trim-index_non_ws <= 0
			if !empty_line && end_critera && !found_end {
				found_end = true
				end_search = index
			}
		}

		if strings.Contains(line, token) && !found_token {
			found_loc = index
			found_token = true
			index_non_ws = len(line) - len(strings.TrimLeft(line, " \t"))
			fmt.Println("index non ws ", index_non_ws)
		}

		fmt.Println("line:", index, line)
	}
	if found_token {
		fmt.Println("Found token on line", found_loc, "and ended on line ", end_search)
	} else {
		fmt.Println("Token not found")
	}

	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		log.Fatal(err)
	}
}
