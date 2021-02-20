package ui

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shurcooL/githubv4"
	"github.com/skanehira/ght/config"
	"github.com/skanehira/ght/domain"
	"github.com/skanehira/ght/github"
	"github.com/skanehira/ght/utils"
)

var IssueUI *SelectUI

func NewIssueUI() {
	opt := func(ui *SelectUI) {
		// initial query
		queries := []string{
			fmt.Sprintf("repo:%s/%s", config.GitHub.Owner, config.GitHub.Repo),
			"state:open",
		}

		IssueFilterUI.SetQuery(strings.Join(queries, " "))

		ui.getList = func(cursor *string) ([]domain.Item, *github.PageInfo) {
			var queries []string
			query := IssueFilterUI.GetQuery()

			if !strings.Contains(query, "is:issue") {
				queries = append(queries, "is:issue")
			}

			for _, q := range strings.Split(query, " ") {
				// execlude Pull request
				if strings.Contains(q, "type:pr") || strings.Contains(q, "is:pr") {
					continue
				}

				queries = append(queries, q)
			}
			query = strings.Join(queries, " ")
			IssueFilterUI.SetQuery(query)

			v := map[string]interface{}{
				"query":  githubv4.String(query),
				"first":  githubv4.Int(30),
				"cursor": (*githubv4.String)(cursor),
			}
			resp, err := github.GetIssue(v)
			if err != nil {
				log.Println(err)
				return nil, nil
			}

			issues := make([]domain.Item, len(resp.Nodes))
			for i, node := range resp.Nodes {
				issues[i] = node.Issue.ToDomain()
			}
			return issues, &resp.PageInfo
		}

		getSelectedIssues := func() []*domain.Issue {
			var issues []*domain.Issue
			if len(IssueUI.selected) == 0 {
				data := IssueUI.GetSelect()
				if data != nil {
					issues = append(issues, data.(*domain.Issue))
				}
			} else {
				for _, item := range IssueUI.selected {
					issues = append(issues, item.(*domain.Issue))
				}
			}
			return issues
		}

		ui.capture = func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Rune() {
			case 'y':
				var urls []string
				for _, issue := range getSelectedIssues() {
					urls = append(urls, issue.URL)
				}

				url := strings.Join(urls, "\n")
				if err := clipboard.WriteAll(url); err != nil {
					log.Println(err)
				}
				IssueUI.ClearSelected()
			case 'o':
				go func() {
					var wg sync.WaitGroup
					for _, issue := range getSelectedIssues() {
						wg.Add(1)
						go func(issue *domain.Issue) {
							defer wg.Done()
							if err := github.ReopenIssue(issue.ID); err != nil {
								log.Println(err)
								return
							}
							issue.State = "OPEN"
						}(issue)
					}
					wg.Wait()
					IssueUI.ClearSelected()
					IssueUI.UpdateView()
				}()
			case 'c':
				go func() {
					var wg sync.WaitGroup
					for _, issue := range getSelectedIssues() {
						wg.Add(1)
						go func(issue *domain.Issue) {
							defer wg.Done()
							if err := github.CloseIssue(issue.ID); err != nil {
								log.Println(err)
								return
							}
							issue.State = "CLOSED"
						}(issue)
					}
					wg.Wait()
					IssueUI.ClearSelected()
					IssueUI.UpdateView()
				}()
			case 'n':
				createIssueForm()
			}
			switch event.Key() {
			case tcell.KeyCtrlO:
				for _, issue := range getSelectedIssues() {
					if err := utils.Open(issue.URL); err != nil {
						log.Println(err)
					}
				}
			}
			return event
		}

		ui.header = []string{
			"",
			"Repo",
			"Number",
			"State",
			"Author",
			"Title",
		}

		ui.hasHeader = len(ui.header) > 0
	}

	ui := NewSelectListUI(UIKindIssue, tcell.ColorBlue, opt)

	ui.SetSelectionChangedFunc(func(row, col int) {
		updateUIRelatedIssue(ui, row)
	})

	IssueUI = ui
}

func createIssueForm() {
	// repo
	var repo string
	input := IssueFilterUI.GetFormItem(0).(*tview.InputField).GetText()
	for _, word := range strings.Split(input, " ") {
		if strings.Contains(word, "repo:") {
			repo = strings.TrimPrefix(word, "repo:")
			break
		}
	}

	if repo == "" {
		return
	}

	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitle("New issue")
	form.SetTitleAlign(tview.AlignLeft)
	inputWidth := 70

	repoInput := tview.NewInputField().SetLabel("Repository").
		SetText(repo).SetLabelWidth(inputWidth).
		SetAcceptanceFunc(func(textToCheck string, lastChar rune) bool {
			return false
		})
	repoInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlN {
			form.SetFocus(2)
		}
		return event
	})

	form.AddFormItem(repoInput)

	s := strings.Split(repoInput.GetText(), "/")
	owner := s[0]
	name := s[1]

	var repoID githubv4.ID
	resp, err := github.GetRepo(map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(name),
	})
	if err != nil {
		log.Println(err)
		return
	}
	repoID = resp.ID

	form.SetFocus(1)

	// autocomplete for assignees, labels, projects, milestones
	autocompleteFunc := func(text string, items []string) []string {
		if text == "" {
			return nil
		}

		words := strings.Split(text, ",")
		word := words[len(words)-1]

		var results []string
		for _, l := range items {
			var isDuplicate bool
			for _, w := range words[:len(words)-1] {
				if w == l {
					isDuplicate = true
					break
				}
			}

			if isDuplicate {
				continue
			}

			if strings.Contains(strings.ToLower(l), strings.ToLower(word)) {
				words = append(words[:len(words)-1], l)
				results = append(results, strings.Join(words, ","))
			}
		}
		return results
	}

	// graphql query variables
	v := map[string]interface{}{
		"owner":  githubv4.String(owner),
		"name":   githubv4.String(name),
		"first":  githubv4.Int(100),
		"cursor": (*githubv4.String)(nil),
	}

	// title
	titleInput := tview.NewInputField().SetLabel("Title").SetLabelWidth(inputWidth)
	form.AddFormItem(titleInput)

	// assignees
	assigneesInput := tview.NewInputField().SetLabel("Assignees").SetLabelWidth(inputWidth)
	form.AddFormItem(assigneesInput)
	userMap := map[string]githubv4.ID{}
	go func() {
		resp, err := github.GetRepoAssignableUsers(v)
		if err != nil {
			log.Println(err)
			return
		}

		if len(resp.Nodes) == 0 {
			UI.app.QueueUpdateDraw(func() {
				form.RemoveFormItem(2)
			})
			return
		}

		var users []string
		for _, u := range resp.Nodes {
			name := string(u.Login)
			userMap[name] = u.ID
			users = append(users, name)
		}
		assigneesInput.SetAutocompleteFunc(func(text string) []string {
			return autocompleteFunc(text, users)
		})
	}()

	// labels
	labelInput := tview.NewInputField().SetLabel("Labels").SetLabelWidth(inputWidth)
	labelMap := map[string]githubv4.ID{}
	form.AddFormItem(labelInput)
	go func() {
		resp, err := github.GetRepoLabels(v)
		if err != nil {
			log.Println(err)
			return
		}

		if len(resp.Nodes) == 0 {
			UI.app.QueueUpdateDraw(func() {
				form.RemoveFormItem(3)
			})
			return
		}

		var labels []string
		for _, l := range resp.Nodes {
			name := string(l.Name)
			labelMap[name] = l.ID
			labels = append(labels, name)
		}
		labelInput.SetAutocompleteFunc(func(text string) []string {
			return autocompleteFunc(text, labels)
		})
	}()

	// projects
	projectInput := tview.NewInputField().SetLabel("Projects").SetLabelWidth(inputWidth)
	projectMap := map[string]githubv4.ID{}
	form.AddFormItem(projectInput)
	go func() {
		resp, err := github.GetRepoProjects(v)
		if err != nil {
			log.Println(err)
			return
		}

		if len(resp.Nodes) == 0 {
			UI.app.QueueUpdateDraw(func() {
				form.RemoveFormItem(4)
			})
			return
		}

		var projects []string
		for _, project := range resp.Nodes {
			name := string(project.Name)
			projectMap[name] = project.ID
			projects = append(projects, name)
		}

		projectInput.SetAutocompleteFunc(func(text string) []string {
			return autocompleteFunc(text, projects)
		})
	}()

	// milestones
	milestoneDropDown := tview.NewDropDown().SetLabel("MileStone").SetLabelWidth(inputWidth)
	var milestoneID *githubv4.ID
	form.AddFormItem(milestoneDropDown)
	go func() {
		resp, err := github.GetRepoMillestones(v)
		if err != nil {
			log.Println(err)
			return
		}

		if len(resp.Nodes) == 0 {
			UI.app.QueueUpdateDraw(func() {
				form.RemoveFormItem(5)
			})
			return
		}

		milestones := map[string]*githubv4.ID{}
		var titles []string
		for _, milestone := range resp.Nodes {
			title := string(milestone.Title)
			milestones[title] = &milestone.ID
			titles = append(titles, title)
		}
		milestoneDropDown.SetOptions(titles, func(text string, index int) {
			milestoneID = milestones[text]
		})
	}()

	var issueBody string
	templateDropDown := tview.NewDropDown().SetLabel("Template").SetLabelWidth(inputWidth)
	go func() {
		v := map[string]interface{}{
			"owner": githubv4.String(owner),
			"name":  githubv4.String(name),
		}
		resp, err := github.GetIssueTemplates(v)
		if err != nil {
			log.Println(err)
			return
		}
		if len(resp) == 0 {
			return
		}

		issueTemplates := map[string]string{}
		var names []string
		for _, te := range resp {
			issueTemplates[string(te.Name)] = string(te.Body)
			names = append(names, string(te.Name))
		}

		templateDropDown.SetOptions(names, func(text string, index int) {
			issueBody = issueTemplates[text]
		})
		UI.app.QueueUpdateDraw(func() {
			form.AddFormItem(templateDropDown)
		})
	}()

	form.AddButton("Edit Body", func() {
		UI.app.Suspend(func() {
			f, err := ioutil.TempFile("", "")
			if err != nil {
				log.Println(err)
				return
			}
			defer os.Remove(f.Name())

			if issueBody != "" {
				if _, err := io.Copy(f, strings.NewReader(issueBody)); err != nil {
					log.Println(err)
				}
			}
			f.Close()

			editor := os.Getenv("EDITOR")
			if editor == "" {
				editor = "vim"
			}
			cmd := exec.Command(editor, f.Name())
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				log.Println(err)
				return
			}

			b, err := ioutil.ReadFile(f.Name())
			if err != nil {
				log.Println(err)
				return
			}

			issueBody = string(b)
		})
	})
	form.AddButton("Create", func() {
		input := githubv4.CreateIssueInput{
			Title:        githubv4.String(titleInput.GetText()),
			RepositoryID: repoID,
		}
		if milestoneID != nil {
			input.MilestoneID = milestoneID
		}

		// get assignee users
		var userIDs []githubv4.ID
		if text := assigneesInput.GetText(); text != "" {
			for _, name := range strings.Split(text, ",") {
				userIDs = append(userIDs, userMap[name])
			}
			input.AssigneeIDs = &userIDs
		}

		// get labels
		var labelIDs []githubv4.ID
		if text := labelInput.GetText(); text != "" {
			for _, name := range strings.Split(text, ",") {
				labelIDs = append(labelIDs, labelMap[name])
			}
			if len(labelIDs) > 0 {
				input.LabelIDs = &labelIDs
			}
		}

		// get projects
		var projectIDs []githubv4.ID
		if text := projectInput.GetText(); text != "" {
			for _, name := range strings.Split(text, ",") {
				projectIDs = append(projectIDs, projectMap[name])
			}
			if len(projectIDs) > 0 {
				input.ProjectIDs = &projectIDs
			}
		}

		body := githubv4.String(issueBody)
		input.Body = &body

		if err := github.CreateIssue(input); err != nil {
			log.Println(err)
		} else {
			UI.pages.RemovePage("form").ShowPage("main")
			UI.app.SetFocus(IssueUI)
			go func() {
				time.Sleep(1 * time.Second)
				IssueUI.GetList()
			}()
		}
	})
	form.AddButton("Cancel", func() {
		UI.pages.RemovePage("form").ShowPage("main")
		UI.app.SetFocus(IssueUI)
	})

	UI.pages.AddAndSwitchToPage("form", UI.Modal(form, 100, 19), true).ShowPage("main")
}

func updateUIRelatedIssue(ui *SelectUI, row int) {
	if row > 0 && row <= len(ui.items) {
		issue := ui.items[row-1].(*domain.Issue)
		IssueViewUI.updateView(issue.Body)

		if len(issue.Comments) > 0 {
			CommentUI.SetList(issue.Comments)
			CommentViewUI.updateView(issue.Comments[0].(*domain.Comment).Body)
		} else {
			CommentUI.ClearView()
			CommentViewUI.Clear()
		}

		if len(issue.Assignees) > 0 {
			AssigneesUI.SetList(issue.Assignees)
		} else {
			AssigneesUI.ClearView()
		}

		if len(issue.Labels) > 0 {
			LabelUI.SetList(issue.Labels)
		} else {
			LabelUI.ClearView()
		}

		if len(issue.MileStone) > 0 {
			MilestoneUI.SetList(issue.MileStone)
		} else {
			MilestoneUI.ClearView()
		}

		if len(issue.Projects) > 0 {
			ProjectUI.SetList(issue.Projects)
		} else {
			ProjectUI.ClearView()
		}
	}
}
