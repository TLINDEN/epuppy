package epub

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"path"
)

// a section in the book
type Section struct {
	File, Title, MediaType string
}

// Book epub book
type Book struct {
	Ncx            Ncx       `json:"ncx"`
	Opf            Opf       `json:"opf"`
	Container      Container `json:"-"`
	Mimetype       string    `json:"-"`
	Content        []Content
	fd             *zip.ReadCloser
	CoverImage     []byte
	CoverFile      string
	CoverMediaType string
	Sections       []Section
	dumpxml        bool
}

// Open open resource file
func (p *Book) Open(n string) (io.ReadCloser, error) {
	return p.open(p.filename(n))
}

// Files list resource files
func (p *Book) Files() []string {
	var fns []string
	for _, f := range p.fd.File {
		fns = append(fns, f.Name)
	}
	return fns
}

// -----------------------------------------------------------------------------
func (p *Book) filename(n string) string {
	return path.Join(path.Dir(p.Container.Rootfile.Path), n)
}

func (p *Book) readXML(n string, v interface{}) error {
	fd, err := p.open(n)
	if err != nil {
		return nil
	}

	defer func() {
		if err := fd.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	dec := xml.NewDecoder(fd)

	if err := dec.Decode(v); err != nil {
		return fmt.Errorf("XML decoder error %w", err)
	}

	return nil
}

func (p *Book) readBytes(n string) ([]byte, error) {
	fd, err := p.open(n)
	if err != nil {
		return nil, nil
	}
	defer func() {
		if err := fd.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	data, err := io.ReadAll(fd)
	if err != nil {
		return nil, fmt.Errorf("failed to read contents from %s %w", n, err)
	}

	return data, nil
}

func (p *Book) open(n string) (io.ReadCloser, error) {
	for _, f := range p.fd.File {
		if f.Name == n {
			return f.Open()
		}
	}

	return nil, fmt.Errorf("file %s not exist", n)
}
