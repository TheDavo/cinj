package cinj

import (
	"log"
	"os"
	"path/filepath"
)

type CinjCommand struct {
	Filepath string
	Args     []string
	FileType Filetype
	SuppArgs []string
}

func (cmd CinjCommand) extractContent() ([]string, error) {
	return []string{}, nil
}

// returnAll simply returns all of the content inside of a file.
// This is called when a cinj command does not have any other arguments.
func (cmd CinjCommand) returnAll() (string, error) {
	content, err := os.ReadFile(cmd.Filepath)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return string(content), nil

}

// fileExtForMarkDown returns a Filetype depending on the extension
// of the file found in the cinj command.
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
