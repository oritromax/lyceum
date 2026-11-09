package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	"github.com/oritromax/lyceum/internal/gitx"
	"github.com/oritromax/lyceum/internal/unit"
)

// prTimeout bounds a gh or glab call.
const prTimeout = 2 * time.Minute

// plainChars need no shell quoting in a printed command.
const plainChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_./:@="

// openPR is close's last step (units.md, "Close" step 5): it pushes branch,
// then opens a pull request for closed unit p, whose plan is planPath, from
// branch into base, or prints how. With an origin remote it pushes branch,
// and only branch, never forced: to its upstream when it has one, else with
// -u to origin (gitx.PushArgs); close is the one command that pushes. If the
// push fails, the unit stays closed, no pull request is opened, and it prints
// the push and the pull request step to run by hand. Pushed, with an origin
// on GitHub and gh on PATH it runs `gh pr create`, and with one on GitLab and
// glab it runs `glab mr create`; otherwise it prints the pull request step
// and the merge commands. With no origin it pushes nothing and prints the
// local merge. The tool's output goes to w. The merge is the owner's, and no
// merge event is written. It returns the next action.
func openPR(w io.Writer, root, branch, base, planPath string, p *unit.Plan) (string, error) {
	merge := fmt.Sprintf("git switch %s && git merge --no-ff %s", base, branch)
	url, err := gitx.RemoteURL(root, "origin")
	if err != nil {
		return "", err
	}
	if url == "" {
		_, _ = fmt.Fprintf(w, "No origin remote, so no pull request. To merge, the owner runs:\n  %s\n", merge)
		return "the owner merges the branch: " + merge, nil
	}

	acs := "none"
	if len(p.ACs) > 0 {
		acs = strings.Join(p.ACs, ", ")
	}
	body := fmt.Sprintf("Plan: %s; ACs: %s", planPath, acs)
	var tool string
	var args []string
	switch {
	case strings.Contains(url, "github.com"):
		tool, args = "gh", []string{"pr", "create", "--base", base, "--head", branch, "--title", p.Title, "--body", body}
	case strings.Contains(url, "gitlab"):
		tool, args = "glab", []string{"mr", "create", "--source-branch", branch, "--target-branch", base, "--title", p.Title, "--description", body, "--yes"}
	}
	if tool != "" {
		if _, err := exec.LookPath(tool); err != nil {
			tool = ""
		}
	}
	pushArgs, err := gitx.PushArgs(root, branch)
	if err != nil {
		return "", err
	}
	push := "git " + strings.Join(pushArgs, " ")
	pushErr := gitx.Push(root, pushArgs)
	if pushErr == nil {
		_, _ = fmt.Fprintf(w, "Pushed %s\n", branch)
	}
	if tool == "" {
		merge += " && git push origin " + base
		if pushErr != nil {
			_, _ = fmt.Fprintf(w, "git push did not push %s (%v); the unit is closed. To push the branch:\n  %s\n", branch, pushErr, push)
		}
		_, _ = fmt.Fprintf(w, "Then open a pull request from %s into %s; or, to merge without one, the owner runs:\n  %s\n", branch, base, merge)
		if pushErr != nil {
			return fmt.Sprintf("%s, then open a pull request from %s into %s; the owner merges it", push, branch, base), nil
		}
		return fmt.Sprintf("open a pull request from %s into %s; the owner merges it", branch, base), nil
	}

	quoted := make([]string, len(args))
	for i, a := range args {
		quoted[i] = a
		if a == "" || strings.IndexFunc(a, func(r rune) bool { return !strings.ContainsRune(plainChars, r) }) >= 0 {
			quoted[i] = shellQuote(a)
		}
	}
	create := tool + " " + strings.Join(quoted, " ")
	if pushErr != nil {
		_, _ = fmt.Fprintf(w, "git push did not push %s (%v); the unit is closed, and no pull request was opened. Push the branch:\n  %s\nthen open the pull request:\n  %s\n", branch, pushErr, push, create)
		return fmt.Sprintf("%s, then run: %s", push, create), nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), prTimeout)
	defer cancel()
	c := exec.CommandContext(ctx, tool, args...)
	c.Dir, c.Stdout, c.Stderr = root, w, w
	if err := c.Run(); err != nil {
		var ee *exec.ExitError
		if !errors.As(err, &ee) && ctx.Err() == nil {
			return "", err
		}
		_, _ = fmt.Fprintf(w, "%s did not open the pull request (%v); the unit is closed. Open it by hand:\n  %s\n", tool, err, create)
		return fmt.Sprintf("open the pull request by hand: %s", create), nil
	}
	_, _ = fmt.Fprintf(w, "Opened a pull request from %s into %s with %s\n", branch, base, tool)
	return "the owner merges the pull request", nil
}
