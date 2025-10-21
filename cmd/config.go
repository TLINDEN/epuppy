/*
Copyright © 2025 Thomas von Dein

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
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
	Version string = `v0.0.7`
	Usage   string = `This is epuppy, a terminal ui ebook viewer.

Usage: epuppy [options] <epub file>

Options:
-D --dark                enable dark mode
-s --store-progress      remember reading position
-n --line-numbers        add line numbers
-c --config <file>       use config <file>
-i --cover-image         display cover image
-t --txt                 dump readable content to STDOUT
-x --xml                 dump source xml to STDOUT
-N --no-color            disable colors (or use $NO_COLOR env var)
-d --debug               enable debugging
-h --help                show help message
-v --version             show program version`
)

type Config struct {
	Showversion     bool         `koanf:"version"`        // -v
	Debug           bool         `koanf:"debug"`          // -d
	StoreProgress   bool         `koanf:"store-progress"` // -s
	Darkmode        bool         `koanf:"dark"`           // -D
	LineNumbers     bool         `koanf:"line-numbers"`   // -n
	Dump            bool         `koanf:"txt"`            // -t
	XML             bool         `koanf:"xml"`            // -x
	NoColor         bool         `koanf:"no-color"`       // -n
	Config          string       `koanf:"config"`         // -c
	ColorDark       ColorSetting `koanf:"colordark"`      // comes from config file only
	ColorLight      ColorSetting `koanf:"colorlight"`     // comes from config file only
	ShowHelp        bool         `koanf:"help"`
	ShowCover       bool         `koanf:"cover-image"` // -i
	Colors          Colors       // generated from user config file or internal defaults, respects dark mode
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
	}

	// parse commandline flags
	flagset.BoolP("version", "v", false, "show program version")
	flagset.BoolP("debug", "d", false, "enable debugging")
	flagset.BoolP("dark", "D", false, "enable dark mode")
	flagset.BoolP("store-progress", "s", false, "store reading progress")
	flagset.BoolP("line-numbers", "n", false, "add line numbers")
	flagset.BoolP("txt", "t", false, "dump readable content to STDOUT")
	flagset.BoolP("xml", "x", false, "dump xml to STDOUT")
	flagset.BoolP("no-color", "N", false, "disable colors")
	flagset.BoolP("cover-image", "i", false, "show cover image")
	flagset.BoolP("help", "h", false, "show help")
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
	} else {
		if !conf.Showversion && !conf.ShowHelp {
			flagset.Usage()
			os.Exit(1)
		}
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

	// disable colors if requested by command line
	if conf.NoColor {
		_ = os.Setenv("NO_COLOR", "1")
	}

	return conf, nil
}

func (c *Config) GetConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "epuppy")
}
