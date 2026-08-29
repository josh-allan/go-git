## go-git

This started as a direct port of @viqueen's [git workspaces](https://github.com/viqueen/git-workspaces):

* `git-recent` -- list most recent git branches that you have locally
* `git-squashed` -- list (and delete) all squashed branches that you still have locally
* `git-merged` -- list (and delete) all merged branches that you still have

With some new commands added in:

* `git-differ` -- a side-by-side diff pager, which also covers `--cached` and `--stat`
* `git-who` -- file ownership and contributor breakdown by repo, directory, or file
* `git-catchup` -- what changed on a branch since you last looked

## Building it

Requires at minimum Go 1.26

```
mise run install
```

## Contributing

```
mise run setup
```
