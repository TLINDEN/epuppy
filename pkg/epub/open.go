package epub

import (
	"archive/zip"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Open open a epub file
func Open(fn string, dumpxml bool) (*Book, error) {
	// to find content
	types := regexp.MustCompile(`application/(xml|html|xhtml|htm)`)

	// cleanup regexes
	deanchor := regexp.MustCompile(`#.*$`)
	cleanext := regexp.MustCompile(`^\.`)

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

	// contains the root path
	err = bk.readXML("META-INF/container.xml", &bk.Container)
	if err != nil {
		return &bk, err
	}

	// contains the OPF data
	err = bk.readXML(bk.Container.Rootfile.Path, &bk.Opf)
	if err != nil {
		return &bk, err
	}

	// look for TOC (might be incomplete, see below!)
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

	// to store our final content sections
	sections := []Section{}

	// count the content items in the raw manifest
	var xmlsmanifest int
	for _, item := range bk.Opf.Manifest {
		if types.MatchString(item.MediaType) {
			xmlsmanifest++
		}
	}

	// we have ncx points from the TOC, try this
	if len(bk.Ncx.Points) > 0 {
		for _, block := range bk.Ncx.Points {
			sect := Section{
				File:  "OEBPS/" + block.Content.Src,
				Title: block.Text,
			}

			srcfile := deanchor.ReplaceAllString(block.Content.Src, "")

			for _, file := range bk.Files() {
				if strings.Contains(file, srcfile) {
					sect.File = file
					sect.MediaType = "application/" + cleanext.ReplaceAllString(filepath.Ext(file), "")
					break
				}
			}

			sections = append(sections, sect)
		}

		if len(sections) < xmlsmanifest {
			// TOC  was incomplete, restart  from scratch but  use the
			// OPF Manifest directly

			sections = []Section{}

			for _, item := range bk.Opf.Manifest {
				if types.MatchString(item.MediaType) {
					sect := Section{
						File:      "OEBPS/" + item.Href,
						MediaType: item.MediaType,
					}

					srcfile := deanchor.ReplaceAllString(item.Href, "")

					for _, file := range bk.Files() {
						if strings.Contains(file, srcfile) {
							sect.File = file
							break
						}
					}

					sections = append(sections, sect)
				}
			}
		}
	} else {
		// no TOC, just pull in the files directly
		for _, file := range bk.Files() {
			sections = append(sections,
				Section{
					File:      file,
					MediaType: "application/" + cleanext.ReplaceAllString(filepath.Ext(file), ""),
				})
		}
	}

	// final sections, store
	bk.Sections = sections

	// to  determine content type, we  could use the MediaType  of the
	// items, but unfortunately we do not have them in every case
	//xmltype := regexp.MustCompile(`(xml|DOCTYPE|xhtml)`)

	// now read in the actual xml contents
	for _, section := range sections {
		content, err := bk.readBytes(section.File)
		if err != nil {
			return &bk, err
		}

		if strings.Contains(section.File, bk.CoverFile) {
			bk.CoverImage = content
		}

		ct := Content{Src: section.File, Title: section.Title}

		//if xmltype.MatchString(string(content)) {
		if types.MatchString(section.MediaType) {
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
