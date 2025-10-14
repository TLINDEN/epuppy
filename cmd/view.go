package cmd

import (
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

	for _, point := range book.Ncx.Points {
		if len(point.Content.Body) > 0 {
			buf.WriteString("### " + point.Content.Title)
			buf.WriteString("\r\n\r\n")
			buf.WriteString(point.Content.Body)
			buf.WriteString("\r\n\r\n\r\n\r\n")
		}
	}

	return Pager(conf, head.String(), buf.String())
}
