package bubbleboxer

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/ansi"
)

const (
	newline = "\n"
	space   = " "
)

// Boxer is a root node of a layout tree and holds all the content of the leaves
type Boxer struct {
	Root       Node
	ContentMap map[string]tea.Model
}

// Node is a node in a layout tree and if it has a address it's a leave
type Node struct {
	Address  string
	Children []Node

	VerticalStacked bool

	width  int
	height int
}

// Init satisfiys the tea.Model interface
func (b Boxer) Init() tea.Cmd { return nil }

// Update satisfys the tea.Model interface
func (b Boxer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO handle empty map, zero area

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return b, tea.Quit
		}
	case tea.WindowSizeMsg:
		b.updateSize(msg.Width, msg.Height)
		return b, nil
	}
	return b, nil
}

// View satisfiys the tea.Model interface
func (b Boxer) View() string {
	if b.Root.width <= 0 || b.Root.height <= 0 {
		return "waiting for size information"
	}
	return strings.Join(b.Root.render(b.ContentMap), newline)
}

func (n *Node) render(leaveContent map[string]tea.Model) []string {
	if n.Address != "" {
		// is leave
		v, ok := leaveContent[n.Address]
		if !ok {
			panic(fmt.Sprintf("address '%s' not found", n.Address))
		}
		return strings.Split(v.View(), newline)
	}

	// is node
	if n.VerticalStacked {
		return n.renderVertical(leaveContent)
	}
	return n.renderHorizontal(leaveContent)
}

func (n *Node) renderVertical(leaveContent map[string]tea.Model) []string {
	if len(n.Children) == 0 {
		panic("no children to render - this node should be a leave or should not exist")
	}
	boxes := make([]string, 0, n.height)
	targetWidth := n.width
	for _, child := range n.Children {
		lines := child.render(leaveContent)
		if len(lines) > child.height {
			panic("content has to much lines")
		}
		// check for to wide lines and because we are on it, pad them to correct width.
		for _, line := range lines {
			lineWidth := ansi.PrintableRuneWidth(line)
			if lineWidth > targetWidth {
				panic("to long line")
			}
			line += strings.Repeat(space, targetWidth-lineWidth)
		}
		boxes = append(boxes, lines...)
		// add more lines to boxes to match the Height of the child-box
		for c := 0; c < child.height-len(lines); c++ {
			boxes = append(boxes, strings.Repeat(space, targetWidth))
		}
	}
	return boxes

}
func (n *Node) renderHorizontal(leaveContent map[string]tea.Model) []string {
	if len(n.Children) == 0 {
		panic("no children to render - this node should be a leave or should not exist")
	}
	//            y  x
	var joinedStr [][]string
	targetHeigth := n.height

	// bring all to same height if they are smaller
	for _, boxer := range n.Children {
		if targetHeigth < boxer.height {
			panic("inconsistent size information: child is bigger than parent")
		}

		lines := boxer.render(leaveContent)

		if len(lines) > targetHeigth {
			panic("content has to much lines")
		}
		if len(lines) < targetHeigth {
			lines = append(lines, make([]string, targetHeigth-len(lines))...)
		}
		joinedStr = append(joinedStr, lines)
	}

	lenght := len(joinedStr)
	// Join the horizontal lines together
	var allStr []string
	// y
	for c := 0; c < targetHeigth; c++ {
		fullLine := make([]string, 0, lenght)
		// x
		for i := 0; i < lenght; i++ {
			boxWidth := n.Children[i].width
			line := joinedStr[i][c]
			lineWidth := ansi.PrintableRuneWidth(line)
			if lineWidth > boxWidth {
				panic("content has to wide lines")
			}
			var pad string
			if lineWidth < boxWidth {
				pad = strings.Repeat(space, boxWidth-lineWidth)
			}
			fullLine = append(fullLine, line, pad)
		}
		allStr = append(allStr, strings.Join(fullLine, ""))
	}
	return allStr

}

func (b *Boxer) updateSize(width, height int) {
	if width <= 0 || height <= 0 {
		panic("wont set area to zero or negative")
	}
	b.Root.updateSize(width, height)
}

// recursive seting of the height and width according to the orientation and the amount of children
func (n *Node) updateSize(width, height int) {
	if width <= 0 || height <= 0 {
		panic("wont set area to zero or negative")
	}
	n.width, n.height = width, height
	if n.Address != "" {
		// is leave
		return
	}
	// is node
	lenght := len(n.Children)
	if lenght == 0 {
		panic("no children to render - this node should be a leave or should not exist")
	}
	for i, c := range n.Children {
		if n.VerticalStacked {
			c.updateSize(width, height/lenght)
		}
		c.updateSize(width/lenght, height)
		n.Children[i] = c
	}
}
