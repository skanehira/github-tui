# github-tui
This is a TUI Client for GitHub.
**Still Under Development**

If you are using Vim, you can use [gh.vim](https://github.com/skanehira/gh.vim) instead.

![](https://i.gyazo.com/d7c8ca82e0aeb947f82c10b08d3eba35.png)

## Features
### Implemented
- Issue
  - list
  - create
  - close
  - open
  - open browser
  - preview
  - edit
- Issue comment
  - list
  - preview
  - delete
  - edit
  - add
  - quote reply

### Still Under Development
- Issue
  - add assignees, labels, projects, milestone
  - remove assignees, labels, projects, milestone
- PR
  - list
  - edit comment
  - add comment
  - delete comment
  - diff
  - create
  - close
  - change base
  - merge
- Github Actions
  - re-run
  - list
  - log
- File tree
  - preview
  - open browser
  - preview
- Project
  - columns
  - open(if type is issue, pr)
  - add
  - remove
  - move
  - open browser
- config
  - set default editor
  - set user keybindings

## Installation

```sh
$ git clone https://github.com/skanehira/github-tui
$ go install ./cmd/ght
```

## Settings
At first, please set personal access [token](https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token) and email in config.yaml.

```yaml
github:
  token: xxxxxxxxxxxxxxxxx
```

The config.yaml path must be in the bellow place.

| OS         | place                                               |
|------------|-----------------------------------------------------|
| Window     | `%AppData%¥ght¥config.yaml`                         |
| Mac        | `$HOME/Library/Application Support/ght/config.yaml` |
| Linux/Unix | `$HOME/.config/ght/config.yaml`                     |

## Usage

```sh
# current repository
$ ght

# specified repository
$ ght owner/repo
```

### Keybindings

| UI       | Keybinding           | Description                      |
|----------|----------------------|----------------------------------|
| Common   | `j`/`down arrow`     | Move down by one row.            |
| Common   | `k`/`up arrow`       | Move up by one row.              |
| Common   | `g`/`home`           | Move to the top.                 |
| Common   | `G`/`end`            | Move to the bottom.              |
| Common   | `Ctrl-F`/`page down` | Move down by one page.           |
| Common   | `Ctrl-B`/`page up`   | Move up by one page.             |
| Common   | `Ctrl-N`             | Move next UI.                    |
| Common   | `Ctrl-P`             | Move previous UI.                |
| Common   | `Ctrl-C`             | Finish app.                      |
| Common   | `Ctrl-G`             | Focus to Issues                  |
| Common   | `Ctrl-T`             | Focus to Filters                 |
| Filters  | `Enter`              | Search with enter query.         |
| Issues   | `h`/`left arrow`     | Move left by one column.         |
| Issues   | `l`/`right arrow`    | Move right by one column.        |
| Issues   | `Ctrl-J`             | Check issue and move down.       |
| Issues   | `Ctrl-K`             | Check issue and move up.         |
| Issues   | `e`                  | Edit and update issue body.      |
| Issues   | `o`                  | Open checked issue.              |
| Issues   | `c`                  | Close checked issue.             |
| Issues   | `Ctrl-O`             | Open checked issue on browser.   |
| Issues   | `/`                  | filter with enter words          |
| Issues   | `n`                  | Create new issue.                |
| Issues   | `f`                  | Fetch more issue.                |
| Comments | `h`/`left arrow`     | Move left by one column.         |
| Comments | `l`/`right arrow`    | Move right by one column.        |
| Comments | `Ctrl-J`             | Check comment and move down.     |
| Comments | `Ctrl-K`             | Check comment and move up.       |
| Comments | `Ctrl-O`             | Open checked comment on browser. |
| Comments | `n`                  | Add new issue comment.           |
| Comments | `e`                  | Edit and update comment body.    |
| Comments | `r`                  | Quote reply comment.             |
| Comments | `/`                  | filter with enter words          |
| Preview  | `/`                  | search with enter words          |
| Preview  | `n`                  | move next word                   |
| Preview  | `N`                  | move previous word               |
| Preview  | `o`                  | change to full screen            |

### Note
When you creating issue, you can specify multiple labels, projects and assignees with `,`.
For instance, when you specify 2 labels then must input `label1,label2`.

![](https://i.gyazo.com/fb665369057c5f096517a24e606e7884.png)

When you edit issue body with `Edit Body` button then `$EDITOR` be used.
If `$EDITOR` is empty or not set, `vim` wll be used.

## Author
skanehira
