package ui

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/skanehira/ght/domain"
	"github.com/skanehira/ght/utils"
)

var ProjectUI *SelectUI

func NewProjectUI() {
	//getList := func(cursor *string) ([]List, github.PageInfo) {
	//	v := map[string]interface{}{
	//		"owner":  githubv4.String(config.GitHub.Owner),
	//		"name":   githubv4.String(config.GitHub.Repo),
	//		"first":  githubv4.Int(100),
	//		"cursor": (*githubv4.String)(cursor),
	//	}
	//	resp, err := github.GetRepoProjects(v)
	//	if err != nil {
	//		return nil, github.PageInfo{}
	//	}

	//	projects := make([]List, len(resp.Nodes))
	//	for i, m := range resp.Nodes {
	//		projects[i] = &Project{
	//			Name: string(m.Name),
	//		}
	//	}

	//	return projects, resp.PageInfo
	//}

	setOpt := func(ui *SelectUI) {
		ui.capture = func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyCtrlO:
				var urls []string
				if len(ProjectUI.selected) == 0 {
					data := ProjectUI.GetSelect()
					if data != nil {
						urls = append(urls, data.(*domain.Project).URL)
					}
				} else {
					for _, s := range ProjectUI.selected {
						urls = append(urls, s.(*domain.Project).URL)
					}
				}

				for _, url := range urls {
					if err := utils.OpenBrowser(url); err != nil {
						log.Println(err)
					}
				}
			}
			return event
		}
	}

	ui := NewSelectListUI(UIKindProject, tcell.ColorLightSalmon, setOpt)
	ProjectUI = ui
}
