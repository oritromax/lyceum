package cli

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/oritromax/lyceum/internal/check"
	"github.com/oritromax/lyceum/internal/config"
	"github.com/oritromax/lyceum/internal/events"
	"github.com/oritromax/lyceum/internal/gate"
	"github.com/oritromax/lyceum/internal/gitx"
	"github.com/oritromax/lyceum/internal/unit"
)

// reopenBug first sends open bug p back from close to cause when a file its
// cause approval hashed changed since (check.CauseDrift), so the owner
// approves cause again; then sends a bug at cause, or still at close, back to fix when a commit
// after its latest verify clear changes a file that is neither a test
// (TESTING.md test_files) nor under lyceum/ (check.Reopened), so verify must
// clear again before approve cause. Each move writes a gate-reopen event with
// the reason, moves the plan's gate (and, to fix, its status to in-progress)
// and commits the plan and the event log. It returns the reasons, "" for a
// move not made.
func reopenBug(root string, p *unit.Plan) (cause, fix string, err error) {
	if p == nil {
		return "", "", nil
	}
	evs, err := events.Read(root)
	if err != nil {
		return "", "", err
	}
	planPath := path.Join(planDir, p.ID+".md")
	drift, err := check.CauseDrift(root, planPath, p, evs)
	if err != nil {
		return "", "", err
	}
	if len(drift) > 0 {
		cause = check.CauseDriftReason(drift)
		if err := reopenTo(root, p, gate.Cause, p.Status, cause); err != nil {
			return "", "", err
		}
		if evs, err = events.Read(root); err != nil {
			return "", "", err
		}
	}
	// At close only code in the branch's own diff counts (AC204).
	var own []string
	if p.Type == string(gate.Bug) && p.Gate == string(gate.Close) && p.Status != unit.StatusClosed {
		cfg, err := config.Read(root)
		if err != nil {
			return "", "", err
		}
		base, err := baseRef(root, cfg.BaseBranch, "reopen the bug", false)
		if err != nil {
			return "", "", err
		}
		if own, err = gitx.BranchCommitted(root, base); err != nil {
			return "", "", err
		}
	}
	fix, err = check.Reopened(root, p, evs, own)
	if err != nil || fix == "" {
		return cause, "", err
	}
	if err := reopenTo(root, p, gate.Fix, unit.StatusInProgress, fix); err != nil {
		return "", "", err
	}
	return cause, fix, nil
}

// reopenTo moves p back to gate g with status, writing a gate-reopen event
// with reason, and commits the plan and the event log.
func reopenTo(root string, p *unit.Plan, g gate.Gate, status, reason string) error {
	// The plan is changed in memory before the event is written, so a plan the
	// CLI cannot rewrite leaves no event behind.
	if err := p.SetCLIField("gate", string(g)); err != nil {
		return err
	}
	if err := p.SetCLIField("status", status); err != nil {
		return err
	}
	if _, err := events.Append(root, events.Event{Type: "gate-reopen", Unit: p.ID,
		Fields: map[string]any{"gate": string(g), "status": status, "reason": reason}}); err != nil {
		return err
	}
	planPath := path.Join(planDir, p.ID+".md")
	if err := unit.WritePlan(filepath.Join(root, filepath.FromSlash(planPath)), p); err != nil {
		return err
	}
	return gitx.Commit(root, fmt.Sprintf("lyceum: reopen %s %s", g, p.ID), planPath, events.Path)
}

// reopenedLine is the line saying unit id went back to fix for reason.
func reopenedLine(id, reason string) string {
	return fmt.Sprintf("Back to %s for %s: %s; verify must clear again before approve cause\n", gate.Fix, id, reason)
}

// causeReopenedLine is the line saying unit id went back to cause for reason.
func causeReopenedLine(id, reason string) string {
	return fmt.Sprintf("Back to %s for %s: %s; the owner approves cause again\n", gate.Cause, id, reason)
}

// reopenedLines are the lines for the moves reopenBug made, in order.
func reopenedLines(id, cause, fix string) string {
	s := ""
	if cause != "" {
		s += causeReopenedLine(id, cause)
	}
	if fix != "" {
		s += reopenedLine(id, fix)
	}
	return s
}
