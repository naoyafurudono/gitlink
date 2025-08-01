package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <file>:<line> or %s <file>:<startLine>-<endLine>\n", os.Args[0], os.Args[0])
		os.Exit(1)
	}

	arg := os.Args[1]

	parts := strings.Split(arg, ":")
	if len(parts) != 2 {
		fmt.Fprintf(os.Stderr, "Invalid format. Use <file>:<line> or <file>:<startLine>-<endLine>\n")
		os.Exit(1)
	}

	filePath := parts[0]
	lineSpec := parts[1]

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting absolute path: %v\n", err)
		os.Exit(1)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "File does not exist: %s\n", filePath)
		os.Exit(1)
	}

	branch, err := getCurrentBranch()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current branch: %v\n", err)
		os.Exit(1)
	}

	remoteURL, err := getRemoteURL()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting remote URL: %v\n", err)
		os.Exit(1)
	}

	repoPath, err := convertToGitHubURL(remoteURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting remote URL: %v\n", err)
		os.Exit(1)
	}

	repoRoot, err := getRepoRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting repository root: %v\n", err)
		os.Exit(1)
	}

	relPath, err := filepath.Rel(repoRoot, absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting relative path: %v\n", err)
		os.Exit(1)
	}

	relPath = filepath.ToSlash(relPath)

	var lineFragment string
	if strings.Contains(lineSpec, "-") {
		rangeParts := strings.Split(lineSpec, "-")
		if len(rangeParts) != 2 {
			fmt.Fprintf(os.Stderr, "Invalid line range format\n")
			os.Exit(1)
		}
		startLine, err1 := strconv.Atoi(rangeParts[0])
		endLine, err2 := strconv.Atoi(rangeParts[1])
		if err1 != nil || err2 != nil {
			fmt.Fprintf(os.Stderr, "Invalid line numbers\n")
			os.Exit(1)
		}
		lineFragment = fmt.Sprintf("#L%d-L%d", startLine, endLine)
	} else {
		line, err := strconv.Atoi(lineSpec)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid line number\n")
			os.Exit(1)
		}
		lineFragment = fmt.Sprintf("#L%d", line)
	}

	url := fmt.Sprintf("%s/blob/%s/%s%s", repoPath, branch, relPath, lineFragment)
	fmt.Println(url)
}

func getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func getRemoteURL() (string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func getRepoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func convertToGitHubURL(remoteURL string) (string, error) {
	url := remoteURL

	// Handle SSH URLs (git@host:path format)
	if strings.HasPrefix(url, "git@") {
		// Extract host and path from git@host:path format
		sshPrefix := "git@"
		remaining := strings.TrimPrefix(url, sshPrefix)
		parts := strings.SplitN(remaining, ":", 2)
		if len(parts) == 2 {
			host := parts[0]
			path := parts[1]
			url = fmt.Sprintf("https://%s/%s", host, path)
		} else {
			return "", fmt.Errorf("invalid SSH URL format: %s", remoteURL)
		}
	} else if strings.HasPrefix(url, "git://") {
		// Convert git:// to https://
		url = strings.Replace(url, "git://", "https://", 1)
	} else if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
		// If it's not a recognized format, return an error
		return "", fmt.Errorf("unsupported remote URL format: %s", remoteURL)
	}

	// Remove .git suffix if present
	url = strings.TrimSuffix(url, ".git")

	return url, nil
}

func readFileLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// Hi!
//
// Enjoy!
