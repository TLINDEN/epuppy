package epub

import (
	"encoding/xml"
	"fmt"
	"strings"
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
	title := Title{}

	err := xml.Unmarshal(content, &title)
	if err != nil {
		if !strings.HasPrefix(err.Error(), "XML syntax error") {
			return fmt.Errorf("XML parser error %w", err)
		}
	}

	c.Title = strings.TrimSpace(title.Content)

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
