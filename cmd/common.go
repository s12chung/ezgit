package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	git "github.com/go-git/go-git/v5"
)

func openRepo() (*git.Repository, *git.Worktree, error) {
	r, err := git.PlainOpen(".")
	if err != nil {
		return nil, nil, err
	}
	w, err := r.Worktree()
	if err != nil {
		return nil, nil, err
	}
	return r, w, nil
}

func askForConfirmation(s, yes string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [%s/n]: ", s, yes)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if strings.ToLower(response) == yes {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

func printStatusAndConfirm(status *git.Status, confirmAsk, confirmText string) bool {
	fmt.Println(status)
	confirm := askForConfirmation(confirmAsk, confirmText)
	if !confirm {
		fmt.Println("Selected 'No', skipping.")
		return false
	}
	return true
}
