package cmd

import (
	"fmt"
	"strings"

	"github.com/tlinden/epuppy/pkg/epub"
)

func View(conf *Config) (int, error) {
	book, err := epub.Open(conf.Document)
	if err != nil {
		return 0, err
	}
	defer book.Close()

	buf := strings.Builder{}
	head := strings.Builder{}

	for _, creator := range book.Opf.Metadata.Creator {
		head.WriteString(creator.Data)
		head.WriteString(" ")
	}

	head.WriteString("- ")

	for _, title := range book.Opf.Metadata.Title {
		head.WriteString(title)
		head.WriteString(" ")
	}

	// FIXME:  since the  switch to  book.Files() in  epub.Open() this
	// returns invalid chapter numbering
	chapter := 1

	for _, content := range book.Content {
		if len(content.Body) > 0 {
			buf.WriteString(conf.Colors.Chapter.
				Render(fmt.Sprintf("Chapter %d: %s", chapter, content.Title)))
			buf.WriteString("\r\n\r\n")
			buf.WriteString(conf.Colors.Body.Render(content.Body))
			buf.WriteString("\r\n\r\n\r\n\r\n")
			chapter++
		}
	}

	return Pager(conf, head.String(), buf.String())
}
