package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	cinj "github.com/TheDavo/cinj/cinj"
)

func main() {
	var cinj cinj.Cinj
	var fp string
	var newname string

	flag.StringVar(&fp, "fp", "", "Filepath of the Markdown file to Cinj")
	flag.StringVar(&newname, "newname", "", "New name for output file, not including extension")

	flag.Parse()
	if fp == "" {
		log.Fatal("No file provided")
		os.Exit(1)
	}

	fmt.Println(fp)
	absFp, err := filepath.Abs(fp)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(absFp)

	_, fileNameNoExt, err := getExtension(absFp)

	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}
	cinj.Filepath = absFp

	if newname == "" {
		cinj.Newname = fileNameNoExt + ".md"
	} else {
		cinj.Newname = filepath.Join(filepath.Dir(absFp), newname+".md")
	}

	err = cinj.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("newname", cinj.Newname)
}

// getExtension is a helper function that returns the correct Cinj-allowed
// extensions or an error
// The allowed extensions are ".cinj.md" or ".cinj"
// This function returns ".cinj.md" or ".cinj" depending on extension
func getExtension(fn string) (string, string, error) {
	fileNoCinj, foundCinj := strings.CutSuffix(fn, ".cinj")
	fileNoCinjMd, foundCinjMd := strings.CutSuffix(fn, ".cinj.md")

	if !foundCinj && !foundCinjMd {
		return "", "", errors.New("No appropriate file extension found")
	}

	if foundCinj {
		return ".cinj", fileNoCinj, nil
	}

	if foundCinjMd {
		return ".cinj.md", fileNoCinjMd, nil
	}

	return "", "", errors.New("No appropriate file extension found")
}
