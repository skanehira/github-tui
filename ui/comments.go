package ui

import (
	"fmt"
	"log"
	"strconv"
	"strings"

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
				deleteComment()
			case 'n':
				item := IssueUI.GetSelect()
				if item == nil {
					return event
				}
				issue := item.(*domain.Issue)

				if err := createComment(item, issue.Body); err != nil {
					UI.Message(err.Error(), func() {
						UI.app.SetFocus(CommentUI)
					})
				}
			case 'e':
				if err := editComment(); err != nil {
					UI.Message(err.Error(), func() {
						UI.app.SetFocus(CommentUI)
					})
				}
			case 'r':
				if err := quoteReply(); err != nil {
					UI.Message(err.Error(), func() {
						UI.app.SetFocus(CommentUI)
					})
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

func quoteReply() error {
	item := CommentUI.GetSelect()
	if item == nil {
		return domain.ErrNotFoundComment
	}

	comment := item.(*domain.Comment)
	lines := strings.Split(comment.Body, "\n")
	for i := range lines {
		lines[i] = fmt.Sprintf("> %s", lines[i])
	}

	body := strings.Join(lines, "\n")

	item = IssueUI.GetSelect()
	if item == nil {
		return domain.ErrNotFoundIssue
	}
	if err := createComment(item, body); err != nil {
		return err
	}

	return nil
}

func createComment(item domain.Item, body string) error {
	if err := editCommentBody(&body); err != nil {
		return err
	}

	input := githubv4.AddCommentInput{
		SubjectID: githubv4.ID(item.Key()),
		Body:      githubv4.String(body),
	}

	if err := github.AddIssueComment(input); err != nil {
		return err
	}

	if err := updateCommentUI(); err != nil {
		return err
	}

	return nil
}

func deleteComment() {
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
}

func editComment() error {
	item := CommentUI.GetSelect()
	if item == nil {
		return domain.ErrNotFoundComment
	}

	comment := item.(*domain.Comment)
	oldBody := comment.Body

	if err := editCommentBody(&comment.Body); err != nil {
		return err
	}

	// if comment body is not changed, do nothing
	if oldBody == comment.Body {
		return nil
	}

	input := githubv4.UpdateIssueCommentInput{
		ID:   githubv4.ID(comment.ID),
		Body: githubv4.String(comment.Body),
	}

	if err := github.UpdateIssueComment(input); err != nil {
		return err
	}

	if err := updateCommentUI(); err != nil {
		return err
	}
	return nil
}

func editCommentBody(body *string) (err error) {
	UI.app.Suspend(func() {
		err = utils.Edit(body)
	})
	if err != nil {
		return
	}

	if *body == "" {
		return domain.ErrCommentBodyIsEmpty
	}
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
		return domain.ErrNotFoundIssue
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
