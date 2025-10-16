package epub

import (
	"testing"
)

func TestEpub(t *testing.T) {
	_, err := open(t, "test.epub")
	if err != nil {
		t.Fatal(err)
	}
}

func open(t *testing.T, f string) (*Book, error) {
	bk, err := Open(f, false)
	if err != nil {
		return nil, err
	}

	t.Logf("files: %+v", bk.Files())
	t.Logf("book: %+v", bk)

	return bk, nil
}
