package cmd

import (
	"fmt"
	"io"
	"log"
)

func Die(err error) int {
	log.Fatal("Error: ", err.Error())

	return 1
}

func Execute(output io.Writer) int {
	conf, err := InitConfig(output)
	if err != nil {
		return Die(err)
	}

	if conf.Showversion {
		_, err := fmt.Fprintf(output, "This is epuppy version %s\n", Version)
		if err != nil {
			return Die(fmt.Errorf("failed to print to output: %s", err))
		}

		return 0
	}

	if conf.StoreProgress {
		progress, err := GetProgress(conf)
		if err == nil {
			conf.InitialProgress = int(progress)
		}
	}

	progress, err := View(conf)
	if err != nil {
		return Die(err)
	}

	if conf.StoreProgress {
		if err := StoreProgress(conf, progress); err != nil {
			return Die(err)
		}
	}

	return 0
}
