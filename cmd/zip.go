package cmd

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var zipCmd = &cobra.Command{
	Use:   "zip [target directory]",
	Short: "Compresses the target directory into a ZIP file and moves it to the parent directory of the target",
	Long: `Compresses the target directory with '_bak' suffix into a ZIP file named after the parent directory of the target.
The resulting ZIP file will be placed in the same directory as the target.`,
	Args: cobra.ExactArgs(1), // このコマンドは正確に1つの引数（ターゲットディレクトリ）を要求します
	Run: func(cmd *cobra.Command, args []string) {
		targetDir := args[0]
		parentDir := filepath.Dir(targetDir)                                    // _bakディレクトリの親ディレクトリを取得
		originalDirName := strings.TrimSuffix(filepath.Base(targetDir), "_bak") // 元のディレクトリ名を取得
		zipFile := filepath.Join(parentDir, originalDirName+".zip")             // ZIPファイルのパスを決定

		fmt.Printf("Compressing directory %s into %s\n", targetDir, zipFile)

		err := compress(targetDir, zipFile)
		if err != nil {
			fmt.Printf("Error compressing directory: %s\n", err)
			return
		}

		fmt.Printf("Directory successfully compressed into ZIP file at %s\n", zipFile)
	},
}

func compress(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	baseSrcDir := filepath.Dir(source)

	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// ソースディレクトリからの相対パスを取得
		relativePath := strings.TrimPrefix(path, baseSrcDir+string(filepath.Separator))

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = relativePath

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
		}
		return err
	})

	if err != nil {
		return err
	}

	return archive.Close()
}

func init() {
	rootCmd.AddCommand(zipCmd)
}
