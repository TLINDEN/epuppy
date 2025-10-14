package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"
	flag "github.com/spf13/pflag"
)

const (
	Version string = `v0.0.1`
	Usage   string = `epuppy [-vd] <epub file>`
)

type Config struct {
	Showversion     bool `koanf:"version"`        // -v
	Debug           bool `koanf:"debug"`          // -d
	StoreProgress   bool `koanf:"store-progress"` // -s
	Document        string
	InitialProgress int // lines
}

func InitConfig(output io.Writer) (*Config, error) {
	var kloader = koanf.New(".")

	// setup custom usage
	flagset := flag.NewFlagSet("config", flag.ContinueOnError)
	flagset.Usage = func() {
		_, err := fmt.Fprintln(output, Usage)
		if err != nil {
			log.Fatalf("failed to print to output: %s", err)
		}
		os.Exit(0)
	}

	// parse commandline flags
	flagset.BoolP("version", "v", false, "show program version")
	flagset.BoolP("debug", "d", false, "enable debugging")
	flagset.BoolP("store-progress", "s", false, "store reading progress")

	if err := flagset.Parse(os.Args[1:]); err != nil {
		return nil, fmt.Errorf("failed to parse program arguments: %w", err)
	}

	// command line setup
	if err := kloader.Load(posflag.Provider(flagset, ".", kloader), nil); err != nil {
		return nil, fmt.Errorf("error loading flags: %w", err)
	}

	// fetch values
	conf := &Config{}
	if err := kloader.Unmarshal("", &conf); err != nil {
		return nil, fmt.Errorf("error unmarshalling: %w", err)
	}

	// arg is the epub file
	if len(flagset.Args()) > 0 {
		conf.Document = flagset.Args()[0]
	}

	return conf, nil
}

func (c *Config) GetConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "epuppy")
}
