package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	boxer "github.com/treilik/bubbleboxer"
)

const (
	upperAddr  = "upper"
	leftAddr   = "left"
	middleAddr = "middle"
	rightAddr  = "right"
	lowerAddr  = "lower"
)

func main() {
	// leaf content creation (models)
	upper := spinnerHolder{spinner.NewModel()}
	left := stringer(leftAddr)
	v := viewport.New(0, 0)
	v.SetContent("foooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo")
	middle := viewPortHolder{v}
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
				1,
				widthOrHeight - 2,
				1,
				// make also sure that the amount of the returned ints match the amount of children:
				// in this case two, but in more complex cases read the amount of the chilren from the len(boxer.Node.Children)
			}
		},
		Children: []boxer.Node{
			stripErr(m.tui.CreateLeaf(upperAddr, upper)),
			{
				Children: []boxer.Node{
					// make sure to encapsulate the models into a leaf with CreateLeaf:
					stripErr(m.tui.CreateLeaf(leftAddr, left)),
					stripErr(m.tui.CreateLeaf(middleAddr, middle)),
					stripErr(m.tui.CreateLeaf(rightAddr, right)),
				},
			},
			stripErr(m.tui.CreateLeaf(lowerAddr, lower)),
		},
	}
	p := tea.NewProgram(m)
	p.EnterAltScreen()
	if err := p.Start(); err != nil {
		fmt.Println(err)
	}
	p.ExitAltScreen()
}

func stripErr(n boxer.Node, _ error) boxer.Node {
	return n
}

type model struct {
	tui boxer.Boxer
}

func (m model) Init() tea.Cmd {
	return spinner.Tick
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.tui.UpdateSize(msg)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.editModel(upperAddr, func(v tea.Model) (tea.Model, error) {
			v, cmd = v.Update(msg)
			return v, nil
		})
		return m, cmd
	}
	return m, nil
}
func (m model) View() string {
	return m.tui.View()
}

func (m *model) editModel(addr string, edit func(tea.Model) (tea.Model, error)) error {
	if edit == nil {
		return fmt.Errorf("no edit function provided")
	}
	v, ok := m.tui.ModelMap[addr]
	if !ok {
		return fmt.Errorf("no model with address '%s' found", addr)
	}
	v, err := edit(v)
	if err != nil {
		return err
	}
	m.tui.ModelMap[addr] = v
	return nil
}

type stringer string

func (s stringer) String() string {
	return string(s)
}

// satisfy the tea.Model interface
func (s stringer) Init() tea.Cmd                           { return nil }
func (s stringer) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return s, nil }
func (s stringer) View() string                            { return s.String() }

type spinnerHolder struct {
	m spinner.Model
}

func (s spinnerHolder) Init() tea.Cmd {
	return s.m.Tick
}
func (s spinnerHolder) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m, cmd := s.m.Update(msg)
	s.m = m
	return s, cmd
}
func (s spinnerHolder) View() string {
	return s.m.View()
}

type viewPortHolder struct {
	m viewport.Model
}

func (v viewPortHolder) Init() tea.Cmd {
	return nil
}
func (v viewPortHolder) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	size, ok := msg.(tea.WindowSizeMsg)
	if ok {
		v.m.Width = size.Width
		v.m.Height = size.Height
		return v, nil
	}
	m, cmd := v.m.Update(msg)
	v.m = m
	return v, cmd
}
func (v viewPortHolder) View() string {
	return v.m.View()
}
