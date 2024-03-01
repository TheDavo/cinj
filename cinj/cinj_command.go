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
