package epub

import (
	"encoding/xml"
	"regexp"
	"strings"
)

var (
	cleantitle    = regexp.MustCompile(`(?s)<head>.*</head>`)
	cleanmarkup   = regexp.MustCompile(`<[^<>]+>`)
	cleanentities = regexp.MustCompile(`&.+;`)
	cleancomments = regexp.MustCompile(`/*.*/`)
	cleanspace    = regexp.MustCompile(`^\s*`)
	cleanh1       = regexp.MustCompile(`<h[1-6].*</h[1-6]>`)
)

// Ncx OPS/toc.ncx
type Ncx struct {
	Points []*NavPoint `xml:"navMap>navPoint" json:"points"`
}

// NavPoint nav point
type NavPoint struct {
	Text    string     `xml:"navLabel>text" json:"text"`
	Content Content    `xml:"content" json:"content"`
	Points  []NavPoint `xml:"navPoint" json:"points"`
}

type Title struct {
	Content string `xml:"head>title"`
}

// Content nav-point content
type Content struct {
	Src   string `xml:"src,attr" json:"src"`
	Empty bool
	Body  string
	Title string
	XML   []byte
}

func (c *Content) String(content []byte) error {
	title := Title{}

	err := xml.Unmarshal(content, &title)
	if err != nil {
		if !strings.HasPrefix(err.Error(), "XML syntax error") {
			return err
		}
	}

	c.Title = title.Content

	txt := cleantitle.ReplaceAllString(string(content), "")
	txt = cleanh1.ReplaceAllString(txt, "")
	txt = cleanmarkup.ReplaceAllString(txt, "")
	txt = cleanentities.ReplaceAllString(txt, " ")
	txt = cleancomments.ReplaceAllString(txt, "")

	txt = strings.TrimSpace(txt)

	c.Body = cleanspace.ReplaceAllString(txt, "")
	c.XML = content

	if len(c.Body) == 0 {
		c.Empty = true
	}

	return nil
}
