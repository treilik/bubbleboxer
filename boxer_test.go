package bubbleboxer

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

type testModel string

func (t testModel) Init() tea.Cmd                       { return nil }
func (t testModel) Update(tea.Msg) (tea.Model, tea.Cmd) { return t, nil }
func (t testModel) View() string                        { return string(t) }

func TestCreateLeaf(t *testing.T) {
	b := Boxer{}
	leaf, err := b.CreateLeaf("test", testModel("test"))
	if err != nil {
		t.Error(err)
	}
	if _, ok := b.ModelMap["test"]; !ok {
		t.Error("after a leaf was created it should be in ModelMap")
	}
	if addr := leaf.GetAddress(); addr != "test" {
		t.Error("after a leaf was created this leaf should have the same address as it was created with")
	}
	if !leaf.IsLeaf() {
		t.Error("a new created Leaf should say of it self that it is a leaf")
	}
}

func TestRenderValidTree(t *testing.T) {
	deferFunc := func() {
		if p := recover(); p != nil {
			t.Errorf("A valid tree should not panic when rendering, but did with: %s", p)
		}
	}
	defer deferFunc()
	b := Boxer{}
	b.LayoutTree = Node{
		Children: []Node{
			{Children: []Node{
				{
					VerticalStacked: true,
					SizeFunc: func(_ Node, widthOrHeight int) []int {
						return []int{
							1,
							widthOrHeight - 1,
						}
					},
					Children: []Node{
						stripErr(b.CreateLeaf("1", testModel("1"))),
						{
							VerticalStacked: true,
							Children: []Node{
								stripErr(b.CreateLeaf("1", testModel("1"))),
								stripErr(b.CreateLeaf("1", testModel("1"))),
							},
						},
					},
				},
				stripErr(b.CreateLeaf("1", testModel("1"))),
			}},
			{
				SizeFunc: func(_ Node, widthOrHeight int) []int {
					return []int{
						1,
						widthOrHeight - 1,
					}
				},
				Children: []Node{
					stripErr(b.CreateLeaf("1", testModel("1"))),
					stripErr(b.CreateLeaf("1", testModel("1"))),
				},
			},
		},
	}
	err := b.UpdateSize(tea.WindowSizeMsg{Width: 17, Height: 22})
	if err != nil {
		t.Error(err)
	}
	b.View()
}
func stripErr(n Node, err error) Node {
	if err != nil {
		panic(err)
	}
	return n
}

func TestMsgHandling(t *testing.T) {
	deferFunc := func() {
		if p := recover(); p != nil {
			t.Errorf("panic while setting up test: '%s'", p)
		}
	}
	defer deferFunc()

	b := Boxer{}
	b.LayoutTree, _ = b.CreateLeaf("test", testModel("test"))

	deferFunc = func() {
		if p := recover(); p != nil {
			t.Errorf("when HandleMsg is true Update should not panic, but did: '%s'", p)
		}
	}
	defer deferFunc()

	b.Update(tea.WindowSizeMsg{Width: 17, Height: 17})
}
