package epub

import (
	"regexp"
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
	Text   string     `xml:"navLabel>text" json:"text"`
	Points []NavPoint `xml:"navPoint" json:"points"`
}

type Title struct {
	Content string `xml:"head>title"`
}
