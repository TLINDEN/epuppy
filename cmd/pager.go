package cmd

// pager setup using bubbletea
// file shamlelessly copied from:
// https://github.com/charmbracelet/bubbletea/tree/main/examples/pager

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"golang.org/x/term"
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()

	viewstyle = lipgloss.NewStyle()
)

const (
	MarginStep = 5
	MinSize    = 40
)

type Meta struct {
	lines           int
	currentline     int
	initialprogress int
}

type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Help     key.Binding
	Quit     key.Binding
	ToggleUI key.Binding
	Pad      key.Binding
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		// every item is one column
		{k.Up, k.Down, k.Left, k.Right},
		{k.Pad}, // fake key, we use it as spacing between columns
		{k.Help, k.Quit, k.ToggleUI},
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

var keys = keyMap{
	Pad: key.NewBinding(
		key.WithKeys("__"),
		key.WithHelp("  ", ""),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("←", "decrease text width"),
	),
	Right: key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("→", "increase text width"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	ToggleUI: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "toggle ui"),
	),
}

type Doc struct {
	content      string
	title        string
	ready        bool
	viewport     viewport.Model
	initialwidth int
	lastwidth    int
	margin       int
	marginMod    bool
	hideUi       bool
	meta         *Meta
	config       *Config

	keys keyMap
	help help.Model
}

func (m Doc) Init() tea.Cmd {
	return nil
}

func (m Doc) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Left):
			if m.lastwidth-(m.margin*2) >= MinSize {
				m.margin += MarginStep
				m.marginMod = true
			}
		case key.Matches(msg, m.keys.Right):
			if m.margin >= MarginStep {
				m.margin -= MarginStep
				m.marginMod = true
			}
		case key.Matches(msg, m.keys.ToggleUI):
			m.hideUi = !m.hideUi
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight

			m.Rewrap()
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
	}

	m.Rewrap()

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// re-calculate word wrapping, add left margin, recalculate line count
func (m *Doc) Rewrap() {
	var content string

	// width has changed, either  because the terminal size changed or
	// because the user added some margin
	if m.lastwidth != m.viewport.Width || m.marginMod {
		content = wordwrap.String(m.content, m.viewport.Width-(m.margin*2))
		m.lastwidth = m.viewport.Width
		m.marginMod = false

		m.viewport.Style = viewstyle.MarginLeft(m.margin)
	}

	// bootstrapping, initialize with default width
	if !m.ready {
		content = wordwrap.String(m.content, m.initialwidth)
	}

	// wrapping has changed, update viewport and line count
	if content != "" {
		m.viewport.SetContent(content)
		m.meta.lines = len(strings.Split(content, "\n"))
	}

	// during bootstrapping: jump to last remembered position, if any
	if !m.ready {
		m.viewport.ScrollDown(m.meta.initialprogress)
	}
}

func (m Doc) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	// update current line for later saving
	m.meta.currentline = int(float64(m.meta.lines) * m.viewport.ScrollPercent())

	var helpView string
	if m.help.ShowAll {
		helpView = "\n" + m.help.View(m.keys)
	}

	if m.hideUi {
		return fmt.Sprintf("%s\n%s", m.viewport.View(), helpView)
	}

	return fmt.Sprintf("%s\n%s\n%s%s", m.headerView(), m.viewport.View(), m.footerView(), helpView)
}

func (m Doc) headerView() string {
	title := m.config.Colors.Title.Render(m.title)
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m Doc) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))

	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Pager(conf *Config, title, message string) (int, error) {
	width := 80
	scrollto := 0

	if term.IsTerminal(int(os.Stdout.Fd())) {
		w, _, err := term.GetSize(0)
		if err == nil {
			width = w
		}
	}

	if conf.StoreProgress {
		scrollto = conf.InitialProgress
	}

	if conf.LineNumbers {
		catn := ""
		for idx, line := range strings.Split(message, "\n") {
			catn += fmt.Sprintf("%4d: %s\n", idx, line)
		}
		message = catn
	}

	meta := Meta{
		initialprogress: scrollto,
		lines:           len(strings.Split(message, "\n")),
	}

	p := tea.NewProgram(
		Doc{
			content:      message,
			title:        title,
			initialwidth: width,
			meta:         &meta,
			config:       conf,
			keys:         keys,
		},
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)

	if _, err := p.Run(); err != nil {
		return 0, fmt.Errorf("could not run pager: %w", err)
	}

	if conf.Debug {
		fmt.Printf("scrollto: %d, last: %d, diff: %d\n",
			scrollto, meta.currentline, scrollto-meta.currentline)
	}

	return meta.currentline, nil
}
