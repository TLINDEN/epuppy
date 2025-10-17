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
	"github.com/charmbracelet/lipgloss"
)

type ColorSetting struct {
	Title   string `koanf:"title"`
	Chapter string `koanf:"chapter"`
	Body    string `koanf:"body"`
}

type Colors struct {
	Title   lipgloss.Style
	Chapter lipgloss.Style
	Body    lipgloss.Style
}

func SetColorconfig(defaultdark, defaultlight ColorSetting, conf *Config) Colors {
	var defaults, user ColorSetting

	switch conf.Darkmode {
	case true:
		defaults = defaultdark
		user = conf.ColorDark
	default:
		defaults = defaultlight
		user = conf.ColorLight
	}

	var colors Colors
	var fg string

	border := lipgloss.RoundedBorder()
	border.Right = "├"
	styletitle := lipgloss.NewStyle().BorderStyle(border).Padding(0, 1)

	if user.Title != "" {
		fg = user.Title
	} else {
		fg = defaults.Title
	}

	colors.Title = styletitle.Foreground(lipgloss.Color(fg))

	if user.Chapter != "" {
		fg = user.Chapter
	} else {
		fg = defaults.Chapter
	}

	colors.Chapter = lipgloss.NewStyle().Foreground(lipgloss.Color(fg))

	if user.Body != "" {
		fg = user.Body
	} else {
		fg = defaults.Body
	}

	colors.Body = lipgloss.NewStyle().Foreground(lipgloss.Color(fg))

	return colors
}
