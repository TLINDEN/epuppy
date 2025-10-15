package epub

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
