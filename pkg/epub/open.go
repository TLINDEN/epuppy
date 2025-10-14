package epub

import (
	"archive/zip"
)

// Open open a epub file
func Open(fn string) (*Book, error) {
	fd, err := zip.OpenReader(fn)
	if err != nil {
		return nil, err
	}

	defer fd.Close()

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
			break
		}
	}

	for _, ncx := range bk.Ncx.Points {
		content, err := bk.readBytes(bk.filename(ncx.Content.Src))
		if err != nil {
			return &bk, err
		}

		if err := ncx.Content.String(content); err != nil {
			return &bk, err
		}
	}

	return &bk, nil
}
