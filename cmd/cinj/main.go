package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Filetype int

const (
	Python Filetype = iota
)

type Cinj struct {
	File   *os.File
	Points []CinjPoint
}

type CinjPoint struct {
	LineNum int
	Command CinjCommand
	Merged  bool
	Type    Filetype
}

type CinjCommand struct {
	File *os.File
	Args []string
}

func main() {
	var filename string
	var newname string

	flag.StringVar(&filename, "fn", "", "Filepath of the Markdown file to Cinj")
	flag.StringVar(&newname, "newname", "", "New name for output file")

	flag.Parse()
	if filename == "" {
		log.Fatal("No file provided")
		os.Exit(1)
	}

	splitFile := strings.Split(filename, ".")
	ext := splitFile[len(splitFile)-1]
	if ext != "md" && ext != "cinj" {
		log.Fatal("Unrecognized or incorrect filetype")
		os.Exit(1)
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(filename)
		log.Fatal(err)
	}

	defer file.Close()
}

func (c *Cinj) FindPoints() []CinjPoint {
	return []CinjPoint{}
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
