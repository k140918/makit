// cmd/root.go
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

const (
	defaultDirPerm  os.FileMode = 0755
	defaultFilePerm os.FileMode = 0644
	timestampLayout             = "200601021504"
)

var (
	// Version は main パッケージから設定されるバージョン文字列（--version で表示）
	Version = "dev"

	mode      string
	timestamp string
	noCreate  bool
	parents   bool
	verbose   bool
	// exitFunc はテスト中に os.Exit をモックできるようにするための変数
	exitFunc = os.Exit
)

// newRootCommand は新しい Cobra コマンドインスタンスを作成し、フラグを初期化します。
// これにより、各テストが独立したコマンドインスタンスを持つことができます。
func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "makit [OPTION] <FILES|DIRS...>",
		Short:   "Create files and directories with optional mode and timestamp",
		Long:    `makit is a CLI tool to create directories and files with optional permissions, timestamps, and parent creation.`,
		Version: Version,
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if verbose {
				fmt.Println("Starting makit operation in verbose mode.")
				fmt.Printf("Arguments: %v\n", args)
				fmt.Printf("Mode: %s, Timestamp: %s, NoCreate: %t, Parents: %t\n", mode, timestamp, noCreate, parents)
			}

			dirPerm, filePerm := defaultDirPerm, defaultFilePerm
			if mode != "" {
				parsed, err := strconv.ParseUint(mode, 8, 32)
				if err != nil || parsed > 0777 {
					if err == nil {
						err = fmt.Errorf("%q is out of range (000-777)", mode)
					}
					fmt.Fprintf(os.Stderr, "invalid mode: %v\n", err)
					exitFunc(1)
					return
				}
				dirPerm = os.FileMode(parsed)
				filePerm = os.FileMode(parsed)
				if verbose {
					fmt.Printf("Parsed permission mode: %o\n", dirPerm)
				}
			}

			var tsTime time.Time
			if timestamp != "" {
				t, err := time.Parse(timestampLayout, timestamp)
				if err != nil {
					fmt.Fprintf(os.Stderr, "invalid timestamp format: %v\n", err)
					exitFunc(1)
					return
				}
				tsTime = t
				if verbose {
					fmt.Printf("Parsed timestamp: %s\n", tsTime.Format("2006-01-02 15:04:05"))
				}
			}

			for _, path := range args {
				if verbose {
					fmt.Printf("Processing path: %s\n", path)
				}

				var perm os.FileMode
				if _, err := os.Stat(path); os.IsNotExist(err) {
					if noCreate {
						fmt.Printf("Skipped (not created): %s\n", path)
						if verbose {
							fmt.Println("Path does not exist and --no-create is set.")
						}
						continue
					}

					if isDirPath(path) {
						perm = dirPerm
						if verbose {
							fmt.Printf("Path '%s' is identified as a directory. Creating...\n", path)
						}
						if err := createDir(path, perm); err != nil {
							fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
							exitFunc(1)
							return
						}
						fmt.Printf("Created directory: %s\n", path)
					} else {
						perm = filePerm
						if verbose {
							fmt.Printf("Path '%s' is identified as a file. Creating...\n", path)
						}
						if err := createFile(path, perm, dirPerm); err != nil {
							fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
							exitFunc(1)
							return
						}
						fmt.Printf("Created file: %s\n", path)
					}
				} else if err != nil {
					fmt.Fprintf(os.Stderr, "Error accessing path: %v\n", err)
					exitFunc(1)
					return
				} else {
					fmt.Printf("Exists: %s\n", path)
					if verbose {
						fmt.Println("Path already exists. Applying mode/timestamp if specified.")
					}
					perm = filePerm
				}

				if mode != "" {
					if verbose {
						fmt.Printf("Applying mode %o to %s\n", perm, path)
					}
					if err := os.Chmod(path, perm); err != nil {
						fmt.Fprintf(os.Stderr, "Warning: failed to set mode on %s: %v\n", path, err)
					}
				}
				if !tsTime.IsZero() {
					if verbose {
						fmt.Printf("Applying timestamp %s to %s\n", tsTime.Format("2006-01-02 15:04:05"), path)
					}
					if err := os.Chtimes(path, tsTime, tsTime); err != nil {
						fmt.Fprintf(os.Stderr, "Warning: failed to set timestamp on %s: %v\n", path, err)
					}
				}
			}
		},
	}

	// テスト時に不要な出力を抑制するため
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	cmd.Flags().StringVarP(&mode, "mode", "m", "", "Set file/directory mode (e.g. 755)")
	cmd.Flags().StringVarP(&timestamp, "date", "d", "", "Set timestamp (e.g. 202504181200)")
	cmd.Flags().BoolVarP(&noCreate, "no-create", "c", false, "Do not create if not exists")
	cmd.Flags().BoolVarP(&parents, "parents", "p", true, "Create parent directories as needed (use --parents=false to disable)")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	return cmd
}

// isDirPath はパスをディレクトリとして扱うかどうかを判定します。
// 末尾が区切り文字の場合、または拡張子がない場合はディレクトリとみなします。
func isDirPath(path string) bool {
	if len(path) > 0 && os.IsPathSeparator(path[len(path)-1]) {
		return true
	}
	return filepath.Ext(path) == ""
}

// createDir はディレクトリを作成します。--parents=false の場合は親ディレクトリを作成しません。
func createDir(path string, perm os.FileMode) error {
	if parents {
		return os.MkdirAll(path, perm)
	}
	return os.Mkdir(filepath.Clean(path), perm)
}

// createFile は空ファイルを作成します。必要に応じて親ディレクトリも作成します。
func createFile(path string, perm, dirPerm os.FileMode) error {
	dir := filepath.Dir(path)
	if dir != "." && parents {
		if verbose {
			fmt.Printf("Ensuring parent directory exists: %s\n", dir)
		}
		if err := os.MkdirAll(dir, dirPerm); err != nil {
			return err
		}
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, perm)
	if err != nil {
		return err
	}
	return f.Close()
}

// Execute はルートコマンドを実行します。
func Execute() {
	cmd := newRootCommand()
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitFunc(1)
	}
}
