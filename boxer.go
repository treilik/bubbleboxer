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

	// HandleMsg controls if update panics or not if its receiving a Msg
	// this is done to make you aware that you should handle all Msg yourself
	// except you know what the Update function does, in this case set HandleMsg to true.
	HandleMsg bool

	// LayoutTree holds the root node and thus the hole LayoutTree
	// Change it as you like as long as every node without children was
	// created with CreateLeaf (to make sure that every leave has a corresponding ModelMap entry)
	// After deleting a Leaf delete the corresponding entry from ModelMap if you care about memory-leaks
	LayoutTree Node

	// ModelMap is a mapping between the Address of a Leaf and the according Model.
	// A valid entry can only be created with CreateLeaf,
	// because entries without a corresponding Node in the LayoutTree are meaningless.
	ModelMap map[string]tea.Model
}

// Node is a node in a layout tree or when created with CreateLeaf its a valid leave of the LayoutTree
type Node struct {
	Children []Node

	// VerticalStacked specifies the orientation of the Children to each other
	VerticalStacked bool

	// SizeFunc specifies the width or height (depending on the orientation) provided to each child.
	// Here by should the sum of the returned int's be the same as the argument 'widthOrHeight'.
	// The length of the returned slice should be the same as the amount of children of the node argument.
	SizeFunc func(node Node, widthOrHeight int) []int

	// noBorder is private because when it changes, the descendants size has to be changed as well
	noBorder bool

	// address is private so that it can only be set if a corresponding entry in Boxer.ModelMap is created (see CreateLeaf)
	address string

	width  int
	height int
}

// SizeError conveys that for at leased one node or leaf in the Layout-tree there was not enough space left
type SizeError error

// Init satisfies the tea.Model interface
func (b Boxer) Init() tea.Cmd { return nil }

// Update panics if HandleMsg is false.
// Otherwise Update reacts to WindowSizeMsg and ctrl+c
func (b Boxer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !b.HandleMsg {
		panic(fmt.Sprintf(`Received Msg: '%s'
but 'HandleMsg' was not set to true.

Either handle all the Msg yourself, so that no Msg reaches this Function
or explicitly set 'HandleMsg' to true, if you know what this does.`, msg))
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return b, tea.Quit
		}
	case tea.WindowSizeMsg:
		b.UpdateSize(msg)
		return b, nil
	}
	return b, nil
}

// View renders the contained tea.Model's according to the LayoutTree
func (b Boxer) View() string {
	if b.LayoutTree.width <= 0 || b.LayoutTree.height <= 0 {
		return "waiting for size information"
	}
	return strings.Join(b.LayoutTree.render(b.ModelMap), NEWLINE)
}

// render recursively renders the layout tree with the models contained in ModelMap
func (n *Node) render(modelMap map[string]tea.Model) []string {
	if n.address != "" {
		// is leaf
		v, ok := modelMap[n.address]
		if !ok {
			panic(fmt.Sprintf("address '%s' not found", n.address))
		}
		return strings.Split(v.View(), NEWLINE)
	}

	// is node
	if n.VerticalStacked {
		return n.renderVertical(modelMap)
	}
	return n.renderHorizontal(modelMap)
}

func (n *Node) renderVertical(modelMap map[string]tea.Model) []string {
	if len(n.Children) == 0 {
		panic("no children to render - this node should be a leaf (see CreateLeaf) or it should not exist")
	}
	boxes := make([]string, 0, n.height)

	targetWidth := n.Children[0].width

	for i, child := range n.Children {
		if child.width != targetWidth {
			panic("inconsistent size information: all children should have the same width when vertical arranged but did not")
		}
		lines := child.render(modelMap)
		if len(lines) > child.height {
			panic("model has to much lines")
		}
		if !n.noBorder && i > 0 {
			lines = append([]string{strings.Repeat(VerticalSeparator, targetWidth)}, lines...)
		}
		// check for to wide lines and because we are on it, pad them to correct width.
		for _, line := range lines {
			lineWidth := ansi.PrintableRuneWidth(line)
			if lineWidth > targetWidth {
				panic("to long line")
			}
			line += strings.Repeat(SPACE, targetWidth-lineWidth)
		}
		boxes = append(boxes, lines...)
		// add more lines to boxes to match the Height of the child-box
		for c := 0; c < child.height-len(lines); c++ {
			boxes = append(boxes, strings.Repeat(SPACE, targetWidth))
		}
	}
	return boxes

}
func (n *Node) renderHorizontal(modelMap map[string]tea.Model) []string {
	if len(n.Children) == 0 {
		panic("no children to render - this node should be a leaf or should not exist")
	}
	//            y  x
	var joinedStr [][]string
	targetHeigth := n.Children[0].height

	// bring all to same height if they are smaller then there own size
	for _, boxer := range n.Children {
		if targetHeigth != boxer.height {
			panic("inconsistent size information: all children should have the same height when horizontal arranged but did not")
		}

		lines := boxer.render(modelMap)

		if len(lines) > targetHeigth {
			panic("model has to much lines")
		}
		if len(lines) < targetHeigth {
			lines = append(lines, make([]string, targetHeigth-len(lines))...)
		}
		joinedStr = append(joinedStr, lines)
	}

	length := len(joinedStr)
	// Join the horizontal lines together
	var allStr []string
	// y
	for c := 0; c < targetHeigth; c++ {
		fullLine := make([]string, 0, length)
		// x
		for i := 0; i < length; i++ {
			boxWidth := n.Children[i].width
			line := joinedStr[i][c]
			lineWidth := ansi.PrintableRuneWidth(line)
			if lineWidth > boxWidth {
				panic("model has to wide lines")
			}
			var pad string
			if lineWidth < boxWidth {
				pad = strings.Repeat(SPACE, boxWidth-lineWidth)
			}
			fullLine = append(fullLine, line+pad)
		}
		var border string
		if !n.noBorder {
			border = HorizontalSeparator
		}

		allStr = append(allStr, strings.Join(fullLine, border))
	}
	return allStr

}

// UpdateSize set the width and height of all Node's
// returns SizeError if the area (width*height) is less or equal to zero for any node or leaf
//
// panics if
//   - a leaf has children
//   - a leaf has a address without a model in the ModelMap (because it was deleted)
//   - a Node (not a leaf) has no Children
//
//   - the SizeFunc returned a slice with different length compared to the size of the Children
//   - the combined space returned by the SizeFunc is greater than the provided size
func (b *Boxer) UpdateSize(size tea.WindowSizeMsg) error {
	return b.LayoutTree.updateSize(size, b.ModelMap)
}

// recursive setting of the height and width according to the orientation and the SizeFunc
// or evenly if no SizeFunc is provided
func (n *Node) updateSize(size tea.WindowSizeMsg, modelMap map[string]tea.Model) error {
	// set size before it may be reduced according to the border
	n.width, n.height = size.Width, size.Height

	// reduce size for children if border is set
	if !n.noBorder {
		length := len(n.Children)
		if length == 0 {
			panic("the border attribute should not be set on a leaf or a node without children")
		}
		// subtract the space which is used by the border between the children
		if n.VerticalStacked {
			size.Height -= length - 1
		} else {
			size.Width -= length - 1
		}
	}

	// check the size after it was reduced
	if size.Width <= 0 || size.Height <= 0 {
		// this returns a error since it is expected that the size might change to to small
		// and return this as a error makes it clear that it is also expected that the calling code has to change the layout
		// according to the size-change or display an alternative message till the size is big enough again.
		return SizeError(fmt.Errorf("not enough space for at least one node or leaf in the Layout-tree"))
	}

	if n.address != "" {
		// is leaf
		if len(n.Children) != 0 {
			panic("a leaf should not have Children")
		}

		v, ok := modelMap[n.address]
		if !ok {
			panic(fmt.Sprintf("no model with address '%s' found", n.address))
		}
		// tell model its size
		v, _ = v.Update(tea.WindowSizeMsg{Width: size.Width, Height: size.Height})
		modelMap[n.address] = v
		return nil
	}

	// is node

	if n.SizeFunc == nil {

		// share space evenly

		length := len(n.Children)
		if length == 0 {
			panic("no children to render - this node should be a leaf or should not exist")
		}
		width := size.Width / length
		height := size.Height

		// hold division remainder (rest)
		restWidth := size.Width % length
		var restHeight int

		if n.VerticalStacked {
			width = size.Width
			height = size.Height / length

			restHeight = size.Height % length
			restWidth = 0
		}

		for i, c := range n.Children {
			var tmpWidth, tmpHeight int
			if restWidth > 0 {
				tmpWidth = 1
				restWidth--
			}
			if restHeight > 0 {
				tmpHeight = 1
				restHeight--
			}

			c.updateSize(
				tea.WindowSizeMsg{
					Width:  width + tmpWidth,
					Height: height + tmpHeight,
				},
				modelMap,
			)
			n.Children[i] = c
		}
		return nil
	}

	// has SizeFunc so split the space according to it
	var sizeList []int
	if n.VerticalStacked {
		sizeList = n.SizeFunc(*n, size.Height)
	} else {
		sizeList = n.SizeFunc(*n, size.Width)
	}
	if len(sizeList) != len(n.Children) {
		panic(fmt.Sprintf("SizeFunc returned %d WindowSizeMsg's but want one for each child and thus: %d", len(sizeList), len(n.Children)))
	}
	var heightSum, widthSum int
	for i, c := range n.Children {
		// set fixed dimension
		s := size

		// change variable dimension according to orientation and the SizeFunc
		if n.VerticalStacked {
			s.Height = sizeList[i]
		} else {
			s.Width = sizeList[i]
		}

		c.updateSize(s, modelMap)
		n.Children[i] = c

		// check sanity
		if n.VerticalStacked {
			heightSum += s.Height
			continue
		}
		widthSum += s.Width
	}

	// the sum of the children size can not be bigger what the parent provided
	if n.VerticalStacked && heightSum > size.Height {
		panic("SizeFunc spread more height than it can")
	}
	if widthSum > size.Width {
		panic("SizeFunc spread more width than it can")
	}
	return nil
}

// CreateLeaf is the only way to create a Node which is treated as a Leaf in the layout-tree.
// CreateLeaf panics when either address is the empty string or the model is nil.
func (b *Boxer) CreateLeaf(address string, model tea.Model) Node {
	if address == "" {
		panic("address should not be empty")
	}
	if model == nil {
		panic("model should not be nil")
	}
	if b.ModelMap == nil {
		b.ModelMap = make(map[string]tea.Model)
	}
	b.ModelMap[address] = model
	return Node{
		address:  address,
		noBorder: true,
	}
}

// IsLeaf returns if the node is a leaf.
func (n *Node) IsLeaf() bool {
	return n.address != ""
}

// GetAddress returns the Address of the Node
// The address of a Node is only settable through CreateLeaf
func (n *Node) GetAddress() string {
	return n.address
}

// GetWidth returns the current with of this node
func (n *Node) GetWidth() int { return n.width }

// GetHeight returns the current with of this node
func (n *Node) GetHeight() int { return n.height }

// CreateNoBorderNode is a constructor for a Node which does not draw a Border around its children.
// The Border attribute is private, because the changing of the attribute has to be accompanied with a change of size
// of all its descendants and is not trivial to facilitate in a save manner.
func CreateNoBorderNode() Node {
	return Node{noBorder: true}
}
