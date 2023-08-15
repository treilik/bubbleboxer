package bubbleboxer

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/ansi"
)

var (
	// NEWLINE is used to separat the lines
	NEWLINE = "\n"
	// SPACE is used to fill up the lines, make sure it is only one column wide and a single character
	SPACE = " "
	// HorizontalSeparator is used to make a visible border between the horizontal arranged children
	// in the layout-tree, make sure it is only one column wide and a single character
	HorizontalSeparator = "│"
	// VerticalSeparator is used to make a visible border between the vertical arranged children
	// in the layout-tree, make sure it is only one column wide and a single character
	VerticalSeparator = "─"
)

// Boxer is a way to render multiple tea.Model's in a specific layout
// according to a LayoutTree.
// The Model's are kept separate from the LayoutTree
// so that changing a Model does not require traversing the LayoutTree.
type Boxer struct {

	// LayoutTree holds the root node and thus the hole LayoutTree
	// Change it as you like as long as every node without children was
	// created with CreateLeaf (to make sure that every leave has a corresponding ModelMap entry)
	// After deleting a Leaf delete the corresponding entry from ModelMap if you care about memory-leaks
	LayoutTree Node

	lastSize tea.WindowSizeMsg
}

// Node is a node in a layout tree or when created with CreateLeaf its a valid leave of the LayoutTree
type Node struct {
	Name string

	// if set to true this node and all its ancestors are ignored
	Hidden bool

	// if model is not nil, this node is treated as a leaf and the Children are ignored
	Model tea.Model

	Children []Node

	// VerticalStacked specifies the orientation of the Children to each other
	VerticalStacked bool

	NoBorder bool

	// SizeFunc specifies the width or height (depending on the orientation) provided to each child.
	// Here by should the sum of the returned int's be the same as the argument 'widthOrHeight'.
	// The length of the returned slice should be the same as the amount of children of the node argument.
	SizeFunc func(node Node, widthOrHeight int) []int
}

func (n Node) String() string {
	return n.Name
}

// SizeError conveys that for at leased one node or leaf in the Layout-tree there was not enough space left
type SizeError error

// NotFoundError convey that the address was not found.
type NotFoundError error

// Init satisfies the tea.Model interface
func (b Boxer) Init() tea.Cmd { return nil }

// Update handles WindowSizeMsg and ctrl+c
func (b Boxer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return b, tea.Quit
		}
	case tea.WindowSizeMsg:
		b.lastSize = msg
		return b, nil
	}
	return b, nil
}

// View renders the contained tea.Model's according to the LayoutTree
func (b Boxer) View() string {
	if b.lastSize.Width <= 0 || b.lastSize.Height <= 0 {
		return "waiting for size information"
	}
	lines, err := b.LayoutTree.render(b.lastSize)
	if err != nil {
		return err.Error()
	}
	return strings.Join(lines, NEWLINE)
}

func even(n Node, widthOrHeight int) []int {
	var visable int
	for _, c := range n.Children {
		if c.Hidden {
			continue
		}
		visable++
	}
	rest := widthOrHeight - (widthOrHeight / visable)
	sizes := make([]int, visable)
	for c := 0; c < visable; c++ {
		sizes[c] = widthOrHeight / visable
		if rest > 0 {
			sizes[c]++
			rest--
		}
	}
	return sizes
}

// render recursively renders the layout tree with the models contained in ModelMap
func (n *Node) render(size tea.WindowSizeMsg) ([]string, error) {
	if n.Model == nil {
		if n.VerticalStacked {
			return n.renderVertical(size)
		}
		return n.renderHorizontal(size)
	}

	m, _ := n.Model.Update(size)

	leaf := strings.Split(m.View(), NEWLINE)
	if len(leaf) > size.Height {
		return leaf, fmt.Errorf("expecting less or equal to %d lines, but the Model with address '%s' has returned to much lines: %d", size.Height, n.Name, len(leaf))
	}

	// pad to correct amount of lines
	leaf = append(leaf, make([]string, size.Height-len(leaf))...)

	for i, line := range leaf {
		lineWidth := ansi.PrintableRuneWidth(line)
		if lineWidth > size.Width {
			return leaf, fmt.Errorf(
				"expecting less or equal to %d character width of all lines, but the Model with address '%s' has returned a to long line with %d characters:%s'%s'",
				size.Width, n.Name, lineWidth, NEWLINE, line,
			)
		}
		// pad to correct width
		leaf[i] = fmt.Sprintf("%s%s", line, strings.Repeat(SPACE, size.Width-lineWidth))
	}
	return leaf, nil
}

func (n *Node) renderVertical(size tea.WindowSizeMsg) ([]string, error) {
	children := make([]Node, 0, len(n.Children))
	for _, c := range n.Children {
		if c.Hidden {
			continue
		}
		children = append(children, c)
	}
	if len(children) == 0 {
		return nil, fmt.Errorf("no children to render - this node should be a leaf (see CreateLeaf) or it should not exist")
	}

	sizes, err := n.sizes(size)
	if err != nil {
		return nil, err
	}

	all := make([]string, len(children))

	var border string
	if !n.NoBorder {
		border = VerticalSeparator
	}
	for i, boxer := range children {
		lines, err := boxer.render(sizes[i])
		if err != nil {
			return lines, wrapError(i, n.VerticalStacked, err)
		}
		all[i] = fmt.Sprintf("%s%s%s", all[i], border, lines)
	}
	return all, nil
}
func (n *Node) renderHorizontal(size tea.WindowSizeMsg) ([]string, error) {
	children := make([]Node, 0, len(n.Children))
	for _, c := range n.Children {
		if c.Hidden {
			continue
		}
		children = append(children, c)
	}
	if len(children) == 0 {
		return nil, fmt.Errorf("no children to render - this node should be a leaf (see CreateLeaf) or it should not exist")
	}

	sizes, err := n.sizes(size)
	if err != nil {
		return nil, err
	}

	all := make([]string, len(children))

	for i, boxer := range children {
		lines, err := boxer.render(sizes[i])
		if err != nil {
			return lines, wrapError(i, n.VerticalStacked, err)
		}
		if i > 0 && !n.NoBorder {
			all = append(all, strings.Repeat(HorizontalSeparator, size.Width))
		}
		all = append(all, lines...)
	}
	return all, nil
}

func (n *Node) sizes(size tea.WindowSizeMsg) ([]tea.WindowSizeMsg, error) {
	var visibleChildren int
	for _, n := range n.Children {
		if n.Hidden {
			continue
		}
		visibleChildren++
	}
	if visibleChildren == 0 {
		return nil, fmt.Errorf("no children to distribute the size to")
	}

	// reduce size for children if border is set
	if !n.NoBorder {
		// subtract the space which is used by the border between the children
		if n.VerticalStacked {
			size.Height -= visibleChildren - 1
		} else {
			size.Width -= visibleChildren - 1
		}
	}

	// check the size after it was reduced
	if size.Width <= 0 || size.Height <= 0 {
		// this returns a error since it is expected that the size might change to to small
		// and return this as a error makes it clear that it is also expected that the calling code has to change the layout
		// according to the size-change or display an alternative message till the size is big enough again.
		return nil, SizeError(fmt.Errorf("not enough space for at least one node or leaf in the Layout-tree"))
	}

	sizeFunc := even
	if n.SizeFunc != nil {
		sizeFunc = n.SizeFunc
	}

	// has SizeFunc so split the space according to it
	var sizeList []int
	if n.VerticalStacked {
		sizeList = sizeFunc(*n, size.Height)
	} else {
		sizeList = sizeFunc(*n, size.Width)
	}

	if len(sizeList) != visibleChildren {
		return nil, fmt.Errorf("SizeFunc returned %d WindowSizeMsg's but want one for each child and thus: %d", len(sizeList), visibleChildren)
	}
	sizes := make([]tea.WindowSizeMsg, visibleChildren)
	for i, s := range sizes {
		if n.VerticalStacked {
			s.Height = sizeList[i]
			s.Width = size.Width
			sizes[i] = s
			continue
		}
		s.Height = size.Height
		s.Width = sizeList[i]
		sizes[i] = s
	}
	return sizes, nil
}

func (b *Boxer) CreateLeaf(address string, model tea.Model) (Node, error) {
	if address == "" {
		return Node{}, fmt.Errorf("empty address given")
	}
	if model == nil {
		return Node{}, fmt.Errorf("no model given")
	}
	err := b.EditNodes(func(n *Node) error {
		if n.Name == address {
			return fmt.Errorf("Address '%s' allready used", address)
		}
		return nil
	})
	if err != nil {
		return Node{}, err
	}
	return Node{
		Name:     address,
		NoBorder: true,
		Model:    model,
	}, nil
}

// EditNodes is called recursivly (after editing) on every node
// if an error occures calling is aborted and the error returned
func (b *Boxer) EditNodes(editFunc func(*Node) error) error {
	return b.LayoutTree.editNodes(editFunc)
}

func (n *Node) editNodes(editFunc func(*Node) error) error {
	err := editFunc(n)
	if err != nil {
		return err
	}
	for _, c := range n.Children {
		err := c.editNodes(editFunc)
		if err != nil {
			return err
		}
	}
	return nil
}

// EditModel edits the models with the given address according to the given funtion if it returns a model and no.
func (b *Boxer) EditModel(address string, editFunc func(tea.Model) (tea.Model, error)) error {
	return b.EditNodes(func(n *Node) error {
		if n.Name != address {
			return nil
		}
		m, err := editFunc(n.Model)
		if err != nil {
			return err
		}
		n.Model = m
		return nil
	})
}

func wrapError(index int, vertical bool, toWrap error) error {
	index++
	layout := "horizontal"
	if vertical {
		layout = "vertical"
	}
	return fmt.Errorf("while rendering the %d child of a %s node a error occured:\n%w", index, layout, toWrap)
}
