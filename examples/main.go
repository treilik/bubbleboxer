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

	m := model{tui: boxer.Boxer{}}
	m.tui.LayoutTree = boxer.Node{
		Children: []boxer.Node{
			{
				VerticalStacked: true,
				Children: []boxer.Node{
					m.tui.CreateLeaf(inAddr, in),
					m.tui.CreateLeaf(mainAddr, main),
				},
			},
			m.tui.CreateLeaf(outAddr, out),
		},
	}
	p := tea.NewProgram(m)
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
		m.tui.UpdateSize(msg)
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
