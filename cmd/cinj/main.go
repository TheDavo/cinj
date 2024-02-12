package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
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
	var filepath string
	var newname string

	flag.StringVar(&filepath, "fp", "", "Filepath of the Markdown file to Cinj")
	flag.StringVar(&newname, "newname", "", "New name for output file")

	flag.Parse()
	if filepath == "" {
		log.Fatal("No file provided")
		os.Exit(1)
	}

	fileext := path.Ext(filepath)

	if fileext != "md" && fileext != "cinj" {
		log.Fatal("Unrecognized or incorrect filetype")
		os.Exit(1)
	}

	if newname != "" {
		newExt := path.Ext(newname)
		if newExt == "" {
			cinj.Newname = newname + ".md"
		} else {
			cinj.Newname = newname
		}
	}

	cinj.Filepath = filepath
	cinj.Run()
}

func (c *Cinj) Run() {
	file, err := os.Open(c.Filepath)
	if err != nil {
		log.Fatal(err)
	}

	newFile, err := os.Create(c.Newname)
	if err != nil {
		log.Fatal(err)
	}

	c.SrcFile = file
	c.DestFile = newFile

	c.cinj()

	defer file.Close()
	defer newFile.Close()
}

func (c *Cinj) cinj() {

	srcScanner := bufio.NewScanner(c.SrcFile)
	destWriter := bufio.NewWriter(c.DestFile)

	for srcScanner.Scan() {
		line := srcScanner.Text()

		if strings.HasPrefix(line, "cinj") {
			command := c.getCinjCommand(line)
			content := c.getContentFromCommand(command)
			contentScanner := bufio.NewScanner(strings.NewReader(content))

			for contentScanner.Scan() {
				contentLine := contentScanner.Text()
				destWriter.WriteString(contentLine)
			}
			srcScanner.Scan()
		}

		destWriter.WriteString(line)

	}
}

func (c Cinj) getCinjCommand(s string) CinjCommand {
	var res CinjCommand
	return res
}

func (c Cinj) getContentFromCommand(cmd CinjCommand) string {
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
