package epub

import (
	"archive/zip"
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// to find content
	types = regexp.MustCompile(`application/(xml|html|xhtml|htm)`)

	// cleanup regexes
	deanchor = regexp.MustCompile(`#.*$`)
	cleanext = regexp.MustCompile(`^\.`)
)

// Open open a epub file and return the filled Book structure
func Open(fn string, dumpxml bool) (*Book, error) {
	bk, err := openFile(fn, dumpxml)
	if err != nil {
		return bk, err
	}

	defer func() {
		if err := bk.fd.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := bk.getManifest(); err != nil {
		return bk, err
	}

	if err := bk.getSections(); err != nil {
		return bk, err
	}

	if err := bk.readSectionContent(); err != nil {
		return bk, err
	}

	return bk, nil
}

// load the epub zip file
func openFile(fn string, dumpxml bool) (*Book, error) {
	fd, err := zip.OpenReader(fn)
	if err != nil {
		return nil, err
	}

	bk := &Book{fd: fd, dumpxml: dumpxml}

	return bk, nil
}

// load the manifest
func (bk *Book) getManifest() error {
	mt, err := bk.readBytes("mimetype")
	if err != nil {
		return err
	}

	bk.Mimetype = string(mt)

	// contains the root path
	err = bk.readXML("META-INF/container.xml", &bk.Container)
	if err != nil {
		return err
	}

	// contains the OPF data
	err = bk.readXML(bk.Container.Rootfile.Path, &bk.Opf)
	if err != nil {
		return err
	}

	// look for TOC (might be incomplete, see below!)
	for _, mf := range bk.Opf.Manifest {
		if mf.ID == bk.Opf.Spine.Toc {
			err = bk.readXML(bk.filename(mf.Href), &bk.Ncx)
			if err != nil {
				return err
			}
		}

		if mf.ID == "cover-image" {
			bk.CoverFile = mf.Href
			bk.CoverMediaType = mf.MediaType
		}
	}

	return nil
}

// extract the readable sections of the epub
func (bk *Book) getSections() error {
	// to store our final content sections
	sections := []Section{}

	// count the content items in the raw manifest
	var manifestcount int
	for _, item := range bk.Opf.Manifest {
		if types.MatchString(item.MediaType) {
			manifestcount++
		}
	}

	// we have ncx points from the TOC, try those
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

		if len(sections) < manifestcount {
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

	// final sections to keep
	bk.Sections = sections

	return nil
}

func (bk *Book) readSectionContent() error {
	// now read in the actual xml contents
	for _, section := range bk.Sections {
		content, err := bk.readBytes(section.File)
		if err != nil {
			return err
		}

		if strings.Contains(section.File, bk.CoverFile) {
			bk.CoverImage = content
		}

		ct := Content{Src: section.File, Title: section.Title}

		if types.MatchString(section.MediaType) {
			if err := ct.String(content); err != nil {
				return err
			}
		}

		if bk.dumpxml {
			fmt.Println(string(ct.XML))
		}

		bk.Content = append(bk.Content, ct)
	}

	return nil
}
