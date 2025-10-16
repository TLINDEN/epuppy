package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	nonprintable = regexp.MustCompile(`[^a-zA-Z0-9\-\._]+`)
	slugify      = regexp.MustCompile(`[/\\]`)
	suffix       = regexp.MustCompile(`\.[a-z]+$`)
)

func StoreProgress(conf *Config, progress int) error {
	cfgpath := conf.GetConfigDir()

	if err := Mkdir(cfgpath); err != nil {
		return fmt.Errorf("failed to mkdir config path %s: %w", cfgpath, err)
	}

	filename := filepath.Join(cfgpath, Slug(conf.Document))

	fd, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open state file %s: %w", filename, err)
	}
	defer fd.Close()

	if err := fd.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate state file %s: %w", filename, err)
	}

	_, err = fd.WriteString(fmt.Sprintf("%d\n", progress))
	if err != nil {
		return fmt.Errorf("failed to write to state file %s: %w", filename, err)
	}

	return nil
}

func GetProgress(conf *Config) (int64, error) {
	cfgpath := conf.GetConfigDir()

	filename := filepath.Join(cfgpath, Slug(conf.Document))

	fd, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		return 0, nil // ignore errors and return no progress
	}
	defer fd.Close()

	scanner := bufio.NewScanner(fd)
	var line string

	for scanner.Scan() {
		line = scanner.Text()
		break
	}

	return strconv.ParseInt(strings.TrimSpace(line), 10, 64)
}

func Mkdir(dir string) error {
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// FIXME: check https://github.com/gosimple/slug
func Slug(input string) string {
	slug := slugify.ReplaceAllString(input, "-")
	slug = suffix.ReplaceAllString(slug, "")
	return nonprintable.ReplaceAllString(slug, "")
}
