package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// gitupdateCmd represents the gitupdate command
var gitupdateCmd = &cobra.Command{
	Use:   "gitupdate [target directory]",
	Short: "Switches the git branch to main or master and pulls the latest changes",
	Long: `Switches the current git branch to 'main' or 'master' (whichever exists) and pulls the latest changes
for the specified directory. This ensures that the directory is up-to-date
before proceeding with other operations.`,
	Args: cobra.ExactArgs(1), // このコマンドは正確に1つの引数（ターゲットディレクトリ）を要求します
	Run: func(cmd *cobra.Command, args []string) {
		targetDir := args[0]
		fmt.Printf("Updating git repository in %s\n", targetDir)

		// ディレクトリに移動
		err := os.Chdir(targetDir)
		if err != nil {
			fmt.Printf("Error changing directory to %s: %s\n", targetDir, err)
			return
		}

		// Check if 'main' branch exists
		if branchExists("main") {
			updateBranch("main")
		} else if branchExists("master") {
			updateBranch("master")
		} else {
			fmt.Println("Neither 'main' nor 'master' branches were found.")
			return
		}
	},
}

func branchExists(branchName string) bool {
	cmd := exec.Command("git", "rev-parse", "--verify", branchName)
	err := cmd.Run()
	return err == nil
}

func updateBranch(branchName string) {
	// git checkout [branchName]
	checkoutCmd := exec.Command("git", "checkout", branchName)
	err := checkoutCmd.Run()
	if err != nil {
		fmt.Printf("Error switching to branch '%s': %s\n", branchName, err)
		return
	}

	// git pull
	pullCmd := exec.Command("git", "pull")
	err = pullCmd.Run()
	if err != nil {
		fmt.Printf("Error pulling latest changes from branch '%s': %s\n", branchName, err)
		return
	}

	fmt.Printf("Git repository successfully updated to branch '%s'.\n", branchName)
}

func init() {
	rootCmd.AddCommand(gitupdateCmd)
	// ここで必要なフラグや設定を追加できます
}
