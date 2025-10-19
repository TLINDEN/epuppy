package epub

import (
	"archive/zip"
	"fmt"
	"log"
	"os"
	"strings"
)

// Open open a epub file
func Open(fn string, dumpxml bool) (*Book, error) {
	fd, err := zip.OpenReader(fn)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := fd.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	bk := Book{fd: fd}
	mt, err := bk.readBytes("mimetype")
	if err != nil {
		return &bk, err
	}

	bk.Mimetype = string(mt)

	err = bk.readXML("META-INF/container.xml", &bk.Container)
	if err != nil {
		return &bk, err
	}

	err = bk.readXML(bk.Container.Rootfile.Path, &bk.Opf)
	if err != nil {
		return &bk, err
	}

	for _, mf := range bk.Opf.Manifest {
		if mf.ID == bk.Opf.Spine.Toc {
			err = bk.readXML(bk.filename(mf.Href), &bk.Ncx)
			if err != nil {
				return &bk, err
			}
		}

		if mf.ID == "cover-image" {
			bk.CoverFile = mf.Href
			bk.CoverMediaType = mf.MediaType
		}
	}

	type section struct {
		file, title string
	}

	sections := []section{}

	if len(bk.Ncx.Points) > 0 {
		for _, block := range bk.Ncx.Points {
			sections = append(sections,
				section{
					file:  "OEBPS/" + block.Content.Src,
					title: block.Text,
				})
		}
	} else {
		for _, file := range bk.Files() {
			sections = append(sections,
				section{
					file: file,
				})
		}
	}

	for _, section := range sections {
		content, err := bk.readBytes(section.file)
		if err != nil {
			return &bk, err
		}

		if strings.Contains(section.file, bk.CoverFile) {
			bk.CoverImage = content
		}

		ct := Content{Src: section.file, Title: section.title}

		if strings.Contains(string(content), "<?xml") {
			if err := ct.String(content); err != nil {
				return &bk, err
			}
		}

		if dumpxml {
			fmt.Println(string(ct.XML))
		}

		bk.Content = append(bk.Content, ct)
	}

	if dumpxml {
		os.Exit(0)
	}

	return &bk, nil
}
