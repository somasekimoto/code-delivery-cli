package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [target directory]",
	Short: "Removes node_modules, .git directories, and files listed in .gitignore within the target directory",
	Long: `Recursively searches the target directory up to 3 levels deep for directories containing a package.json file.
For each found directory, removes the node_modules and .git directories, and any files or directories listed in the .gitignore file.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetDir := args[0]
		fmt.Printf("Searching and cleaning directories in %s\n", targetDir)
		err := searchAndClean(targetDir, 0)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		fmt.Println("Cleaning process completed successfully.")
	},
}

func searchAndClean(path string, level int) error {
	// レベル3より深い階層は走査しない
	if level > 3 {
		return nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		currentPath := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			// package.jsonが含まれるディレクトリを見つけたら、特定のディレクトリとファイルを削除
			if containsPackageJSON(currentPath) {
				err = cleanDirectory(currentPath)
				if err != nil {
					return err
				}
			} else {
				// サブディレクトリを再帰的に走査
				err = searchAndClean(currentPath, level+1)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func containsPackageJSON(path string) bool {
	_, err := os.Stat(filepath.Join(path, "package.json"))
	return !os.IsNotExist(err)
}

func cleanDirectory(path string) error {
	// node_modules と .git ディレクトリを削除
	directoriesToRemove := []string{"node_modules", ".git"}
	for _, dir := range directoriesToRemove {
		dirPath := filepath.Join(path, dir)
		err := os.RemoveAll(dirPath)
		if err != nil {
			return fmt.Errorf("failed to remove %s: %w", dirPath, err)
		}
	}

	// .gitignoreにリストされているファイルやディレクトリを削除
	gitignorePath := filepath.Join(path, ".gitignore")
	gitignoreContents, err := os.ReadFile(gitignorePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read .gitignore: %w", err)
	}
	for _, line := range strings.Split(string(gitignoreContents), "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			itemPath := filepath.Join(path, line)
			err := os.RemoveAll(itemPath)
			if err != nil {
				return fmt.Errorf("failed to remove %s: %w", itemPath, err)
			}
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
