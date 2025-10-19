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
	cleanmarkup   = regexp.MustCompile(`<[^<>]+>`)
)

// Content nav-point content
type Content struct {
	Src   string `xml:"src,attr" json:"src"`
	Empty bool
	Body  string
	Title string
	XML   []byte
}

// parse XML, look for title and <p>.*</p> stuff
func (c *Content) String(content []byte) error {
	doc, err := xmlquery.Parse(
		strings.NewReader(
			cleanentitles.ReplaceAllString(string(content), " ")))
	if err != nil {
		return err
	}

	if c.Title == "" {
		// extract the title
		for _, item := range xmlquery.Find(doc, "//title") {
			c.Title = strings.TrimSpace(item.InnerText())
		}
	}

	// extract all  paragraphs, ignore any formatting  and re-fill the
	// paragraph,  that is, we  replace all newlines inside  with one
	// space.
	txt := strings.Builder{}
	var have_p bool
	for _, item := range xmlquery.Find(doc, "//p") {
		if !empty.MatchString(item.InnerText()) {
			have_p = true
			txt.WriteString(newlines.ReplaceAllString(item.InnerText(), " ") + "\n\n")
		}
	}

	if !have_p {
		// try  <div></div>, which some  ebooks use, so get  all divs,
		// remove markup and paragraphify the parts
		for _, item := range xmlquery.Find(doc, "//div") {
			if !empty.MatchString(item.InnerText()) {
				cleaned := cleanmarkup.ReplaceAllString(item.InnerText(), "")
				txt.WriteString(newlines.ReplaceAllString(cleaned, " ") + "\n\n")
			}
		}
	}

	c.Body = strings.TrimSpace(txt.String())
	c.XML = content

	if len(c.Body) == 0 {
		c.Empty = true
	}

	return nil
}
