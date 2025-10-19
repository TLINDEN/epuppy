package epub

import (
	"regexp"
	"strings"

	"github.com/antchfx/xmlquery"
)

var (
	cleanentitles = regexp.MustCompile(`&[a-z]+;`)
	empty         = regexp.MustCompile(`(?s)^[\s ]*$`)
	newlines      = regexp.MustCompile(`[\r\n]+`)
)

// Content nav-point content
type Content struct {
	Src   string `xml:"src,attr" json:"src"`
	Empty bool
	Body  string
	Title string
	XML   []byte
}

func (c *Content) String(content []byte) error {
	// parse XML, look for title and <p>.*</p> stuff
	doc, err := xmlquery.Parse(
		strings.NewReader(
			cleanentitles.ReplaceAllString(string(content), " ")))
	if err != nil {
		return err
	}

	// extract the title
	for _, item := range xmlquery.Find(doc, "//title") {
		c.Title = strings.TrimSpace(item.InnerText())
	}

	// extract all  paragraphs, ignore any formatting  and re-fill the
	// paragraph,  that is, we  replaces all newlines inside  with one
	// space.
	txt := strings.Builder{}
	for _, item := range xmlquery.Find(doc, "//p") {
		if !empty.MatchString(item.InnerText()) {
			txt.WriteString(newlines.ReplaceAllString(item.InnerText(), " ") + "\n\n")
		}
	}

	c.Body = strings.TrimSpace(txt.String())
	c.XML = content

	if len(c.Body) == 0 {
		c.Empty = true
	}

	return nil
}
