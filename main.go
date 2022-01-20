package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/treilik/bubblesgum/list"
)

func main() {
	content := make(map[string]tea.Model)
	content["in"] = list.NewModel()
	content["main"] = list.NewModel()
	content["out"] = list.NewModel()
	p := tea.NewProgram(
		Boxer{
			ContentMap: content,
			Root: Node{
				Children: []Node{
					{
						Address: "in",
					},
					{
						Address: "main",
					},
					{
						Address: "out",
					},
				},
			},
		},
	)
	p.Start()
}
