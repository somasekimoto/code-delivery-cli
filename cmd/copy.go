package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// copyCmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy [source directory]",
	Short: "Copies the source directory recursively and creates a backup with _bak suffix",
	Long: `Copies the specified source directory recursively and creates a backup of it.
The backup directory will have the same name as the source but with a '_bak' suffix.`,
	Args: cobra.ExactArgs(1), // このコマンドは正確に1つの引数（ソースディレクトリ）を要求します
	Run: func(cmd *cobra.Command, args []string) {
		sourceDir := args[0]
		backupDir := sourceDir + "_bak"

		// バックアップディレクトリが既に存在するか確認
		if _, err := os.Stat(backupDir); !os.IsNotExist(err) {
			fmt.Printf("Backup directory already exists: %s\n", backupDir)
			return
		}

		// ディレクトリを再帰的にコピー
		err := copyDir(sourceDir, backupDir)
		if err != nil {
			fmt.Printf("Error copying directory: %s\n", err)
			return
		}

		fmt.Printf("Directory successfully copied to: %s\n", backupDir)
	},
}

// copyDir recursively copies a directory tree, attempting to preserve permissions.
func copyDir(src string, dst string) error {
	// ソースディレクトリの内容を取得
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// ディレクトリを作成
	err = os.MkdirAll(dst, 0755)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		fileInfo, err := os.Stat(srcPath)
		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			// サブディレクトリの場合は再帰的にコピー
			err = copyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			// ファイルの場合はコピー
			err = copyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// copyFile copies the contents of the file named src to the file named by dst.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func init() {
	rootCmd.AddCommand(copyCmd)
	// ここで必要なフラグや設定を追加できます
}
