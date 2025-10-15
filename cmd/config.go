package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/alecthomas/repr"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"
	flag "github.com/spf13/pflag"
)

const (
	Version string = `v0.0.2`
	Usage   string = `epuppy [-vd] <epub file>`
)

type Config struct {
	Showversion   bool         `koanf:"version"`        // -v
	Debug         bool         `koanf:"debug"`          // -d
	StoreProgress bool         `koanf:"store-progress"` // -s
	Darkmode      bool         `koanf:"dark"`           // -D
	Config        string       `koanf:"config"`         // -c
	ColorDark     ColorSetting `koanf:"colordark"`      // comes from config file only
	ColorLight    ColorSetting `koanf:"colorlight"`     // comes from config file only

	Colors          Colors // generated from user config file or internal defaults, respects dark mode
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
	flagset.BoolP("dark", "D", false, "enable dark mode")
	flagset.BoolP("store-progress", "s", false, "store reading progress")
	flagset.StringP("config", "c", "", "read config from file")

	if err := flagset.Parse(os.Args[1:]); err != nil {
		return nil, fmt.Errorf("failed to parse program arguments: %w", err)
	}

	// generate a  list of config files to try  to load, including the
	// one provided via -c, if any
	var configfiles []string

	configfile, _ := flagset.GetString("config")
	home, _ := os.UserHomeDir()

	if configfile != "" {
		configfiles = []string{configfile}
	} else {
		configfiles = []string{
			"/etc/epuppy.toml", "/usr/local/etc/epuppy.toml", // unix variants
			filepath.Join(home, ".config", "epuppy", "config.toml"),
		}
	}

	// Load the config file[s]
	for _, cfgfile := range configfiles {
		if path, err := os.Stat(cfgfile); !os.IsNotExist(err) {
			if !path.IsDir() {
				if err := kloader.Load(file.Provider(cfgfile), toml.Parser()); err != nil {
					return nil, fmt.Errorf("error loading config file: %w", err)
				}
			}
		} // else: we ignore the file if it doesn't exists
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

	if conf.Debug {
		repr.Println(conf)
	}

	// setup color config
	conf.Colors = SetColorconfig(
		ColorSetting{ // Dark
			Title:   "#ff4500",
			Chapter: "#ff4500",
			Body:    "#cdb79e",
		},
		ColorSetting{ // Light
			Title:   "#ff0000",
			Chapter: "#8b0000",
			Body:    "#696969",
		},
		conf)

	return conf, nil
}

func (c *Config) GetConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "epuppy")
}
