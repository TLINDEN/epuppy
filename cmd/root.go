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
	"os"
)

func Die(err error) int {
	fmt.Fprintln(os.Stderr, "Error: ", err.Error())

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

	if conf.ShowHelp {
		_, err := fmt.Fprintln(output, Usage)
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
