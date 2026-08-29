## Go-git

This is a trimmed-down Go port of @viqueen's [git workspaces](https://github.com/viqueen/git-workspaces):

* git-recent -- list most recent git branches
* git-squashed -- list (and delete) all squashed branches
* git-merged -- list (and delete) all merged branches

It also includes a new command:

* git-differ -- a side-by-side diff pager, which also covers --cached and --stat

## Building it

Requires at minimum Go 1.26

```
mise run install
```

`mise run install` installs to GOBIN.

## Contributing

```
mise run setup 
```

This installs both goimports and configures the relevant git hooks.

