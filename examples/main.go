package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	boxer "github.com/treilik/bubbleboxer"
	"github.com/treilik/bubblesgum/list"
)

const (
	inAddr   = "in"
	mainAddr = "main"
	outAddr  = "out"
)

func main() {
	in := list.NewModel()
	main := list.NewModel()
	out := list.NewModel()

	in.AddItems([]fmt.Stringer{stringer(inAddr)})
	main.AddItems([]fmt.Stringer{stringer(mainAddr)})
	out.AddItems([]fmt.Stringer{stringer(outAddr)})

	content := make(map[string]tea.Model)
	content[inAddr] = in
	content[mainAddr] = main
	content[outAddr] = out
	p := tea.NewProgram(
		model{
			tui: boxer.Boxer{
				ContentMap: content,
				Root: boxer.Node{
					Children: []boxer.Node{
						{
							Children: []boxer.Node{
								{
									Address: inAddr,
								},
								{
									Address: mainAddr,
								},
							},
						},
						{
							Address: outAddr,
						},
					},
				},
			},
		},
	)
	p.Start()
}

type model struct {
	tui boxer.Boxer
}

func (m model) Init() tea.Cmd { return nil }
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.tui.UpdateSize(msg.Width, msg.Height)
	}
	return m, nil
}
func (m model) View() string {
	return m.tui.View()
}

type stringer string

func (s stringer) String() string {
	return string(s)
}
