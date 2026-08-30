## go-git

This started as a direct port of @viqueen's [git workspaces](https://github.com/viqueen/git-workspaces):

* `git-recent` -- list most recent git branches that you have locally
* <img width="3018" height="444" alt="image" src="https://github.com/user-attachments/assets/2f21a24c-e7f3-44dc-9930-10727bc2471c" />

* `git-squashed` -- list (and delete) all squashed branches that you still have locally
* <img width="3012" height="392" alt="image" src="https://github.com/user-attachments/assets/e3184e61-7d4e-4861-b254-d53052dd7da7" />

* `git-merged` -- list (and delete) all merged branches that you still have locally
* <img width="3020" height="346" alt="image" src="https://github.com/user-attachments/assets/c75af468-e584-4988-b565-08b5bc329960" />


With some new commands added in:

* `git-differ` -- a side-by-side diff pager, which also covers `--cached` and `--stat`
* <img width="3024" height="764" alt="image" src="https://github.com/user-attachments/assets/1f69521f-d08e-4fd2-ba2c-24cedc0794a7" />

* `git-who` -- file ownership and contributor breakdown by repo, directory, or file
* <img width="3024" height="458" alt="image" src="https://github.com/user-attachments/assets/5ced803f-2f7b-408f-a68c-3277924288b1" />

* `git-catchup` -- what changed on a branch since you last looked
* <img width="3024" height="472" alt="image" src="https://github.com/user-attachments/assets/d4c1eb55-0e56-4120-af9f-fa82e24a91a0" />


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
