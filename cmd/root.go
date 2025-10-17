/*
Copyright © 2025 Ruslan Khalikov rycln1@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var (
	source, target string
	rootCmd        = &cobra.Command{
		Use:   "filer",
		Short: "filer - Interactive file sorting REPL: keep or delete files one-by-one",
		Long: `Usage: filer [-s SOURCE] [-t TARGET]

	Interactive file sorting REPL

	Options:
	-s, --source DIR  Source directory (default: current)
	-t, --target DIR  Target directory for kept files (default: keep in place)

	Commands: [K]eep, [D]elete, [Q]uit
	Example: filer -s ~/Downloads -t ~/Keep`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		Run: func(cmd *cobra.Command, args []string) {
			entries, err := os.ReadDir(source)
			if err != nil {
				fmt.Println(err)
				return
			}

			if target != "" {
				err := os.MkdirAll(target, 0755)
				if err != nil {
					fmt.Println(err)
					return
				}
			}

			scanner := bufio.NewScanner(os.Stdin)
			for _, entry := range entries {
				fmt.Printf("Name: %s\n[K]eep, [D]elete, [Q]uit?", entry.Name())
				if !scanner.Scan() {
					break
				}

				input := strings.TrimSpace(scanner.Text())
				if input == "" {
					continue
				}

				switch input {
				case "q", "Q":
					return
				case "k", "K":
					if target == "" {
						continue
					}
					err := moveFileSafe(source+"/"+entry.Name(), target+"/"+entry.Name())
					if err != nil {
						fmt.Println(err)
						return
					}
				case "d", "D":
					err := os.Remove(source + "/" + entry.Name())
					if err != nil {
						fmt.Println(err)
						return
					}
				}
			}
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.filer.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringVarP(&source, "source", "s", "./", "Source directory (default: current)")
	rootCmd.Flags().StringVarP(&target, "target", "t", "", "Target directory for kept files (default: keep in place)")
}

func moveFileSafe(sourcePath, destPath string) error {
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", sourcePath)
	}

	err := os.Rename(sourcePath, destPath)
	if err == nil {
		return nil
	}

	return copyAndRemove(sourcePath, destPath)
}

func copyAndRemove(sourcePath, destPath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		os.Remove(destPath)
		return err
	}

	sourceInfo, _ := sourceFile.Stat()
	destInfo, _ := destFile.Stat()

	if sourceInfo.Size() != destInfo.Size() {
		os.Remove(destPath)
		return fmt.Errorf("the file sizes do not match")
	}

	return os.Remove(sourcePath)
}
