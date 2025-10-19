/*
Copyright © 2025 Thomas von Dein

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/repr"
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

	return Pager(&Ebook{
		Config: conf,
		Title:  conf.Document,
		Body:   string(data),
	})
}

func ViewEpub(conf *Config) (int, error) {
	book, err := epub.Open(conf.Document, conf.XML)
	if err != nil {
		return 0, err
	}

	if conf.Debug {
		repr.Println(book.Opf)
		repr.Println(book.Ncx)
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

	return Pager(&Ebook{
		Config:    conf,
		Title:     head.String(),
		Body:      buf.String(),
		Cover:     book.CoverImage,
		MediaType: book.CoverMediaType,
	})
}

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
