## go-git

This started as a direct port of @viqueen's [git workspaces](https://github.com/viqueen/git-workspaces):

* `git-recent` -- list most recent git branches that you have locally
* `git-squashed` -- list (and delete) all squashed branches that you still have locally
* `git-merged` -- list (and delete) all merged branches that you still have

With some new commands added in:

* `git-differ` -- a side-by-side diff pager, which also covers `--cached` and `--stat`
* `git-who` -- file ownership and contributor breakdown by repo, directory, or file
* `git-catchup` -- what changed on a branch since you last looked
* <img width="3014" height="944" alt="image" src="https://github.com/user-attachments/assets/483e66b5-bb94-45aa-9c6e-ab7a89b0cb93" />

## Installation

Requires Go 1.26+

```
go install ./cmd/...
```

This installs all binaries to your `GOBIN` (typically `~/go/bin`). Make sure that's on your `PATH`.

With [mise](https://mise.jdx.dev):

```
mise run install
```

## Contributing

```
go install golang.org/x/tools/cmd/goimports@latest
git config core.hooksPath .githooks
```

Or with mise:

```
mise run setup
```
