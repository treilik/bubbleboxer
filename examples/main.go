package main

import (
	tea "github.com/charmbracelet/bubbletea"
	boxer "github.com/treilik/bubbleboxer"
	"github.com/treilik/bubblesgum/list"
)

func main() {
	content := make(map[string]tea.Model)
	content["in"] = list.NewModel()
	content["main"] = list.NewModel()
	content["out"] = list.NewModel()
	p := tea.NewProgram(
		boxer.Boxer{
			ContentMap: content,
			Root: boxer.Node{
				Children: []boxer.Node{
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
