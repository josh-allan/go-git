package root

import (
	"context"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:   "go-git",
	Short: "go git yourself some git hygiene",
}

func Execute() {
	if err := fang.Execute(
		context.Background(),
		rootCmd,
		fang.WithVersion(version),
		fang.WithNotifySignal(os.Interrupt),
	); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(recentCmd())
	rootCmd.AddCommand(mergedCmd())
	rootCmd.AddCommand(squashedCmd())
}
