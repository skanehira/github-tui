package ui

import (
	"errors"
	"log"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/shurcooL/githubv4"
	"github.com/skanehira/ght/domain"
	"github.com/skanehira/ght/github"
	"github.com/skanehira/ght/utils"
	"golang.org/x/sync/errgroup"
)

var CommentUI *SelectUI

func NewCommentUI() {
	setOpt := func(ui *SelectUI) {
		ui.capture = func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Rune() {
			case 'd':
				UI.Confirm("Do you want to delete comments?", "Yes", func() error {
					comments := getSelectedComments()
					if len(comments) == 0 {
						return nil
					}

					var eg errgroup.Group
					for _, comment := range comments {
						id := comment.ID
						eg.Go(func() error {
							return github.DeleteIssueComment(id)
						})
					}

					// When all the processing is completed this error be returned
					// because if some of delete action be success, need to update view
					deleteErr := eg.Wait()
					if deleteErr != nil {
						log.Println(deleteErr)
					}

					if err := updateCommentUI(); err != nil {
						return err
					}

					return deleteErr
				}, func() {
					UI.app.SetFocus(CommentUI)
				})

			case 'n':
				focus := func() {
					UI.app.SetFocus(CommentUI)
				}
				item := IssueUI.GetSelect()
				if item == nil {
					UI.Message("not found issue", focus)
					return event
				}

				var body string

				if err := editCommentBody(&body); err != nil {
					UI.Message(err.Error(), focus)
					return event
				}

				input := githubv4.AddCommentInput{
					SubjectID: githubv4.ID(item.Key()),
					Body:      githubv4.String(body),
				}

				if err := github.AddIssueComment(input); err != nil {
					UI.Message(err.Error(), focus)
					return event
				}

				if err := updateCommentUI(); err != nil {
					UI.Message(err.Error(), focus)
					return event
				}
			case 'e':
				item := ui.GetSelect()
				if item == nil {
					return event
				}

				focus := func() {
					UI.app.SetFocus(CommentUI)
				}

				comment := item.(*domain.Comment)
				oldBody := comment.Body

				if err := editCommentBody(&comment.Body); err != nil {
					UI.Message(err.Error(), focus)
					return event
				}

				// if comment body is not changed, do nothing
				if oldBody == comment.Body {
					return event
				}

				input := githubv4.UpdateIssueCommentInput{
					ID:   githubv4.ID(comment.ID),
					Body: githubv4.String(comment.Body),
				}

				if err := github.UpdateIssueComment(input); err != nil {
					UI.Message(err.Error(), focus)
					return event
				}

				if err := updateCommentUI(); err != nil {
					UI.Message(err.Error(), focus)
					return event
				}
			}

			switch event.Key() {
			case tcell.KeyCtrlO:
				for _, comment := range getSelectedComments() {
					if err := utils.Open(comment.URL); err != nil {
						log.Println(err)
					}
				}
				CommentUI.ClearSelected()
				CommentUI.UpdateView()
			}

			return event
		}

		ui.header = []string{
			"",
			"Author",
			"UpdatedAt",
		}
		ui.hasHeader = len(ui.header) > 0
	}

	ui := NewSelectListUI(UIKindComment, tcell.ColorYellow, setOpt)

	ui.SetSelectionChangedFunc(func(row, col int) {
		if row > 0 {
			CommentViewUI.updateView(ui.items[row-1].(*domain.Comment).Body)
		}
	})

	CommentUI = ui
}

func editCommentBody(body *string) (err error) {
	UI.app.Suspend(func() {
		err = utils.Edit(body)
	})
	return
}

func getSelectedComments() []*domain.Comment {
	var comments []*domain.Comment
	if len(CommentUI.selected) == 0 {
		data := CommentUI.GetSelect()
		comments = append(comments, data.(*domain.Comment))
	} else {
		for _, item := range CommentUI.selected {
			comments = append(comments, item.(*domain.Comment))
		}
	}
	return comments
}

func updateCommentUI() error {
	item := IssueUI.GetSelect()
	if item == nil {
		return errors.New("not found issue")
	}
	oldIssue := item.(*domain.Issue)

	number, err := strconv.Atoi(oldIssue.Number)
	if err != nil {
		return err
	}
	m := map[string]interface{}{
		"owner":  githubv4.String(oldIssue.RepoOwner),
		"name":   githubv4.String(oldIssue.Repo),
		"number": githubv4.Int(number),
	}

	issue, err := github.GetIssue(m)
	if err != nil {
		return err
	}

	newIssue := issue.ToDomain()
	IssueUI.UpdateItem(newIssue)

	if len(newIssue.Comments) > 0 {
		CommentUI.SetList(newIssue.Comments)
		CommentViewUI.updateView(newIssue.Comments[0].(*domain.Comment).Body)
	} else {
		CommentUI.ClearView()
		CommentViewUI.Clear()
	}
	return nil
}
