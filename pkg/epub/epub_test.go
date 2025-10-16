package epub

import (
	"log"
	"testing"
)

func TestEpub(t *testing.T) {
	bk, err := open(t, "test.epub")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := bk.Close(); err != nil {
			log.Fatal(err)
		}
	}()

}

func open(t *testing.T, f string) (*Book, error) {
	bk, err := Open(f)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := bk.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	t.Logf("files: %+v", bk.Files())
	t.Logf("book: %+v", bk)

	return bk, nil
}
