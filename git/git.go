package git

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/fabric8io/gitcontroller/util"
	"strings"
)

// GitClone clones a git repo for the given git url at the given path location
func GitClone(url string, path string) error {
	util.Info("Cloning git repository ")
	util.Success(url)
	util.Info(" at ")
	util.Success(path)
	util.Info("\n")

	cmd := exec.Command("git", "clone", url, path)
	return waitForCommand(cmd)
}

// GitPull performs a git pull
func GitPull(path string) error {
	err := os.Chdir(path)
	if err != nil {
		return err
	}
	cmd := exec.Command("git", "pull")
	return waitForCommandIgnoreOutput(cmd)
}

// GitLatestCommitSince returns the latest commit id in the same branch for the local clone
// of a git repository
func GitLatestCommitSince(path string, ref string) (string, error) {
	err := os.Chdir(path)
	if err != nil {
		return "", err
	}
	branch, err := gitBranchOfRef(ref)
	if err != nil {
		return "", err
	}
	if len(branch) <= 0 {
		branch = "master"
	}
	out, err := exec.Command("git", "log", branch, "-n", "1", "--format=oneline").Output()
	if err != nil {
		return "", err
	}
	return firstWord(string(out)), nil
}

func gitBranchOfRef(ref string) (string, error) {
	out, err := exec.Command("git", "branch", "--contains", ref).Output()
	if err != nil {
		return "", err
	}
	text := string(out)
	if strings.HasPrefix(text, "*") {
		text = text[1:]
	}
	text = strings.TrimSpace(text)
	return firstWord(text), nil
}

func firstWord(text string) string {
	array := strings.Split(text, " ")
	if len(array) < 1 {
		return text
	}
	return array[0]
}

func waitForCommand(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return waitForCommandIgnoreOutput(cmd)
}

func waitForCommandIgnoreOutput(cmd *exec.Cmd) error {
	var waitStatus syscall.WaitStatus
	if err := cmd.Run(); err != nil {
		printErr(err)
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			printStatus(waitStatus.ExitStatus())
		}
		return err
	} else {
		waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
		printStatus(waitStatus.ExitStatus())
		return nil
	}
}

func printStatus(exitStatus int) {
	if exitStatus != 0 {
		util.Error(fmt.Sprintf("%d", exitStatus))
	}
}

func printErr(err error) {
	if err != nil {
		util.Errorf("%s\n", err.Error())
	}
}
