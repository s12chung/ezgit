package cmd

import (
	"bytes"
	_ "embed" //nolint:golint
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"
)

func newStartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start a new branch for a new ep",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStart()
		},
	}
}

var errSkip = fmt.Errorf("skipping")

func runStart() error {
	if err := startCheckout(); err != nil {
		if err == errSkip {
			return nil
		}
		return err
	}
	if err := startGenerateFiles(); err != nil {
		return err
	}
	return nil
}

func startCheckout() error {
	r, w, err := openRepo()
	if err != nil {
		return err
	}
	status, err := w.Status()
	if err != nil {
		return err
	}
	if len(status) != 0 {
		if !printStatusAndConfirm(&status, "Local changes exist, DISCARD them?", "discard") {
			return errSkip
		}
	}

	if err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName("main"),
		Force:  true,
	}); err != nil {
		return err
	}

	if err = w.Pull(&git.PullOptions{RemoteName: "origin"}); err != nil && err != git.NoErrAlreadyUpToDate {
		return nil
	}
	if err = cleanAllBranches(r); err != nil {
		return nil
	}

	newEpisodeNumber, err := startGetLastEpisodeNumber()
	if err != nil {
		return err
	}

	if err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName("ep" + newEpisodeNumber),
		Create: true,
	}); err != nil {
		return err
	}

	return nil
}

func cleanAllBranches(r *git.Repository) error {
	iter, err := r.Branches()
	if err != nil {
		return err
	}

	var branches []string
	err = iter.ForEach(func(ref *plumbing.Reference) error {
		branches = append(branches, ref.Name().String())
		return nil
	})
	if err != nil {
		return err
	}

	for _, branch := range branches {
		if branch == "refs/heads/main" {
			continue
		}
		if err = r.Storer.RemoveReference(plumbing.ReferenceName(branch)); err != nil {
			return err
		}
	}
	return nil
}

//go:embed embed/description.yaml.tmpl
var descriptionTemplate string

//go:embed embed/tweet.txt
var tweetTemplate []byte

const descriptionTemplateName = "description"

func startGenerateFiles() error {
	newEpisodeNumber, err := startGetLastEpisodeNumber()
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile("./tweets/"+newEpisodeNumber+".txt", tweetTemplate, 0600); err != nil {
		return err
	}

	t, err := template.New(descriptionTemplateName).Parse(descriptionTemplate)
	if err != nil {
		return err
	}
	var buffer bytes.Buffer
	err = t.ExecuteTemplate(&buffer, descriptionTemplateName, map[string]string{
		"episodeNumber": newEpisodeNumber,
	})
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile("./descriptions/"+newEpisodeNumber+".yaml", buffer.Bytes(), 0600); err != nil {
		return err
	}
	return nil
}

func startGetLastEpisodeNumber() (string, error) {
	var files []string
	err := filepath.WalkDir(filepath.Join(".", "tweets"), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".txt" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	episodeNumbers := make([]int, len(files))
	for i, file := range files {
		file = filepath.Base(file)
		episodeNumber, err := strconv.ParseInt(strings.TrimSuffix(file, filepath.Ext(file)), 10, 32)
		if err != nil {
			return "", err
		}
		episodeNumbers[i] = int(episodeNumber)
	}
	sort.Ints(episodeNumbers)
	return strconv.Itoa(episodeNumbers[len(episodeNumbers)-1] + 1), nil
}
