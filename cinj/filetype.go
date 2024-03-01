package cinj

type Filetype string

func (ft Filetype) String() string {
	return string(ft)
}

const (
	Python     Filetype = "python"
	Javascript          = "javascript"
	Markdown            = "md"
	Text                = ""
	Plain               = ""
)
