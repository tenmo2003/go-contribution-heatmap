package git

import (
	"contribution-heatmap/internal/utils"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetUserEmailFromGitConfig() (string, error) {
	cmd := exec.Command("git", "config", "user.email")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

func getGitDir(path string) (string, error) {
	gitDir := filepath.Join(path, ".git")
	exists, err := utils.PathExists(gitDir)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", fmt.Errorf("not a git repository")
	}
	return gitDir, nil
}

func GetReposGitDirs(root string) ([]string, error) {
	repos := []string{}
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if path == root {
			return nil
		}
		if err != nil {
			fmt.Println(err)
			return err
		}
		if d.IsDir() {
			gitDir, err := getGitDir(path)
			if err != nil {
				if err.Error() == "not a git repository" {
					return filepath.SkipDir
				}
				return err
			}
			repos = append(repos, gitDir)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return repos, nil
}

func GetCommitCountByDate(author string, from string, to string, repos []string) map[string]int {
	commitCountByDate := map[string]int{}

	for _, repo := range repos {
		cmd := exec.Command(
			"git",
			"--git-dir="+repo,
			"log",
			"--pretty=format:%ad",
			"--reverse",
			"--date=format:%Y-%m-%d",
			"--author="+author,
			"--since="+from,
			"--until="+to,
		)

		out, err := cmd.Output()
		if err != nil {
			continue
		}

		outputStr := strings.TrimSpace(string(out))
		if outputStr == "" {
			continue
		}
		lines := strings.SplitSeq(outputStr, "\n")

		for line := range lines {
			parts := strings.Split(line, " ")
			date := parts[0]

			commitCountByDate[date]++
		}
	}
	return commitCountByDate
}
