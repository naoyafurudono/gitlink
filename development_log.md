# Development Log

## 2025-08-26

### Enhancement: Use Commit Hash Instead of Branch Name
- Changed URL generation to use commit hash instead of branch name for more stable links
- Replaced `getCurrentBranch()` function with `getCurrentCommitHash()`
- Now uses `git rev-parse HEAD` to get the current commit SHA
- URLs now point to specific commits rather than branches, ensuring links remain valid even after branch changes
- Example URL format: `https://github.com/user/repo/blob/{commit_hash}/file#L{line}`

## 2025-08-01

### Bug Fix: git.pepabo.com Support
- Fixed error "unsupported remote URL format: https://git.pepabo.com/hosting/gulliver.git"
- Added support for git.pepabo.com remote URLs in `convertToGitHubURL` function:
  - SSH format: `git@git.pepabo.com:`
  - HTTPS format: `https://git.pepabo.com/`
  - Git protocol format: `git://git.pepabo.com/`
- Updated main.go:135-139 to handle pepabo.com domains alongside existing GitHub support

## 2025-07-22

### Initial Implementation
- Created `gitlink` CLI tool based on readme.md specifications
- Implemented core functionality:
  - Parse command-line arguments in format `file:line` or `file:startLine-endLine`
  - Get current git branch using `git rev-parse --abbrev-ref HEAD`
  - Get remote origin URL using `git config --get remote.origin.url`
  - Convert various git URL formats to GitHub web URLs
  - Generate GitHub blob URLs with line number fragments
- Created Go module with go.mod file
- Tool successfully generates GitHub links for specified code ranges