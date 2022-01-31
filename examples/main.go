package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	boxer "github.com/treilik/bubbleboxer"
)

const (
	leftAddr   = "left"
	middleAddr = "middle"
	rightAddr  = "right"
	lowerAddr  = "lower"
)

func main() {
	// leaf content creation (models)
	left := stringer(leftAddr)
	middle := stringer(middleAddr)
	right := stringer(rightAddr)

	lower := stringer(fmt.Sprintf("%s: use ctrl+c to quit", lowerAddr))

	// layout-tree defintion
	m := model{tui: boxer.Boxer{}}
	m.tui.LayoutTree = boxer.Node{
		// orientation
		VerticalStacked: true,
		// spacing
		SizeFunc: func(_ boxer.Node, widthOrHeight int) []int {
			return []int{
				// since this node is vertical stacked return the height partioning since the width stays for all children fixed
				widthOrHeight - 1,
				1,
				// make also sure that the amount of the returned ints match the amount of children:
				// in this case two, but in more complex cases read the amount of the chilren from the len(boxer.Node.Children)
			}
		},
		Children: []boxer.Node{
			{
				Children: []boxer.Node{
					// make sure to encapsulate the models into a leaf with CreateLeaf:
					m.tui.CreateLeaf(leftAddr, left),
					m.tui.CreateLeaf(middleAddr, middle),
					m.tui.CreateLeaf(rightAddr, right),
				},
			},
			m.tui.CreateLeaf(lowerAddr, lower),
		},
	}
	p := tea.NewProgram(m)
	p.EnterAltScreen()
	if err := p.Start(); err != nil {
		fmt.Println(err)
	}
	p.ExitAltScreen()
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

// satisfy the tea.Model interface
func (s stringer) Init() tea.Cmd                           { return nil }
func (s stringer) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return s, nil }
func (s stringer) View() string                            { return s.String() }
