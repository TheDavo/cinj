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

var cinjDescription = `Cinj is a command line tool that expands on the markdown syntax to make
code report generation easier. 

By modifying code block syntax in a markdown file,
automating code-heavy report generation is now an easier process, especially when paired
with tools such as 'pandoc' and 'weasyprint' for PDF generation.`

func main() {
	var cinj cinj.Cinj
	var newname string

	flag.StringVar(
		&newname,
		"newname",
		"",
		"New name for output file, not including extension,\n\tfor example --newname new_report_name",
	)

	flag.Usage = func() {
		w := flag.CommandLine.Output()

		fmt.Fprintf(w, "Usage of %s\n", os.Args[0])
		fmt.Fprintf(w, "cinj [options] filepath\n")
		fmt.Fprintf(w, "%s\n\n", cinjDescription)
		fmt.Fprintln(w, "Options and flags:")
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	absFp, err := filepath.Abs(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

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
	fmt.Println("Newly Cinj'd filename:", cinj.Newname)
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
