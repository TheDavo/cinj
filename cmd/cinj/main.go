package main

import (
	"flag"
	"fmt"
	cinj "github.com/TheDavo/cinj/cinj"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var cinj cinj.Cinj
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
		temp1, found1 := strings.CutSuffix(absFp, ".cinj")
		temp2, found2 := strings.CutSuffix(absFp, ".md")
		temp3, found3 := strings.CutSuffix(absFp, ".cinj.md")
		if found3 {
			cinj.Newname = temp3 + ".md"
		} else if found2 {
			cinj.Newname = temp2 + ".cinj.md"
		} else if found1 {
			cinj.Newname = temp1 + ".md"
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

	err = cinj.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("newname", cinj.Newname)
}
