package epub

import (
	"archive/zip"
	"fmt"
	"log"
	"os"
	"strings"
)

// Open open a epub file
func Open(fn string) (*Book, error) {
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

			break
		}
	}

	dump := os.Getenv("DUMPXML")

	for _, file := range bk.Files() {
		content, err := bk.readBytes(file)
		if err != nil {
			return &bk, err
		}

		ct := Content{Src: file}
		if strings.Contains(string(content), "<?xml") {
			if err := ct.String(content); err != nil {
				return &bk, err
			}
		}

		bk.Content = append(bk.Content, ct)

		if dump != "" {
			fmt.Println(string(ct.XML))
		}
	}

	if dump != "" {
		os.Exit(0)
	}

	return &bk, nil
}
