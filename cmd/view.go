package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tlinden/epuppy/pkg/epub"
)

func View(conf *Config) (int, error) {
	switch filepath.Ext(conf.Document) {
	case ".epub":
		return ViewEpub(conf)
	default:
		return ViewText(conf)
	}
}

func ViewText(conf *Config) (int, error) {
	data, err := os.ReadFile(conf.Document)
	if err != nil {
		return 0, err
	}

	if conf.Dump {
		return fmt.Println(string(data))
	}

	return Pager(conf, conf.Document, string(data))
}

func ViewEpub(conf *Config) (int, error) {
	book, err := epub.Open(conf.Document, conf.XML)
	if err != nil {
		return 0, err
	}

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

	fetchByContent(conf, &buf, book)

	if conf.Dump {
		return fmt.Println(buf.String())
	}

	return Pager(conf, head.String(), buf.String())
}

// FIXME:  since the  switch to  book.Files() in  epub.Open() this
// returns invalid chapter numbering
func fetchByContent(conf *Config, buf *strings.Builder, book *epub.Book) bool {
	chapter := 1
	var gotbody bool

	for _, content := range book.Content {
		if len(content.Body) > 0 {
			if content.Title != "" {
				buf.WriteString(conf.Colors.Chapter.
					Render(fmt.Sprintf("──────┤ %s ├──────", content.Title)))

			}
			buf.WriteString("\r\n\r\n")

			if conf.Dump {
				// avoid excess whitespaces when printing to stdout
				buf.WriteString(content.Body)
			} else {
				buf.WriteString(conf.Colors.Body.Render(content.Body))
			}

			buf.WriteString("\r\n\r\n\r\n\r\n")
			chapter++

			gotbody = true
		}
	}

	return gotbody
}
