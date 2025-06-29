package app

import (
	"fmt"

	"github.com/Lew-Lew/sherlock/style"
)

type Result struct {
	URL   string
	Label string
	Issue string
	Err   error
}

func NewResult(url, label, issue string) Result {
	return Result{
		URL:   url,
		Label: label,
		Issue: issue,
	}
}

func NewErrorResult(url string, err error) Result {
	return Result{
		URL: url,
		Err: err,
	}
}

func (rl Result) String() string {
	if rl.Err != nil {
		return fmt.Sprintf("%v\n%v",
			style.IssueKey.Render(rl.URL),
			style.IssueAlert.Render(fmt.Sprintf("Error: %v", rl.Err)),
		)
	}
	return fmt.Sprintf("%v %v\n%v",
		style.IssueKey.Render(rl.URL),
		style.IssueAlert.Render(rl.Issue),
		style.IssueInfo.Render(rl.Label),
	)
}
