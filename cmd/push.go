package cmd

import (
	"fmt"

	git "github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

func newPushCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "push",
		Short: "Push local changs to Github",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPush()
		},
	}
}

func runPush() error {
	r, w, err := openRepo()
	if err != nil {
		return err
	}
	status, err := w.Status()
	if err != nil {
		return err
	}

	if err := pushCommit(r, w, status); err != nil {
		if err == errNoChanges {
			return nil
		}
		return err
	}

	return pushShowPRURL(r)
}

var errNoChanges = fmt.Errorf("no changes")

func pushCommit(r *git.Repository, w *git.Worktree, status git.Status) error {
	if len(status) == 0 {
		fmt.Println("No changes to push")
		return errNoChanges
	}

	for file := range status {
		_, err := w.Add(file)
		if err != nil {
			return err
		}
	}

	commit, err := w.Commit("changes", &git.CommitOptions{})
	if err != nil {
		return err
	}
	if _, err = r.CommitObject(commit); err != nil {
		return err
	}

	if err = r.Push(&git.PushOptions{}); err != nil {
		return err
	}

	return err
}

func pushShowPRURL(r *git.Repository) error {
	head, err := r.Head()
	if err != nil {
		return err
	}
	fmt.Printf("Push Successful, to create a Pull Request: https://github.com/s12chung/trailheadspodcast/compare/%s?expand=1\n", head.Name().Short())
	return err
}
