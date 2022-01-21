package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	boxer "github.com/treilik/bubbleboxer"
	"github.com/treilik/bubblesgum/list"
)

const (
	upperAddr  = "upper"
	leftAddr   = "left"
	middleAddr = "middle"
	rightAddr  = "right"
)

func main() {
	upper := list.NewModel()
	left := list.NewModel()
	middle := list.NewModel()
	right := list.NewModel()

	upper.AddItems([]fmt.Stringer{stringer("use 'ctrl+c' or 'q' to quit")})
	left.AddItems([]fmt.Stringer{stringer(leftAddr)})
	middle.AddItems([]fmt.Stringer{
		stringer(middleAddr),
	})
	right.AddItems([]fmt.Stringer{stringer(rightAddr)})

	m := model{tui: boxer.Boxer{}}

	m.tui.LayoutTree = boxer.Node{
		VerticalStacked: true,
		SizeFunc: func(_ boxer.Node, msg tea.WindowSizeMsg) []tea.WindowSizeMsg {
			return []tea.WindowSizeMsg{
				{Height: msg.Height - 1, Width: msg.Width},
				{Height: 1, Width: msg.Width},
			}
		},
		Children: []boxer.Node{
			{
				Children: []boxer.Node{
					m.tui.CreateLeaf(leftAddr, left),
					m.tui.CreateLeaf(middleAddr, middle),
					m.tui.CreateLeaf(rightAddr, right),
				},
			},
			m.tui.CreateLeaf(upperAddr, upper),
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
