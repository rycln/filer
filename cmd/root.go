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
	"regexp"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// rootCmd represents the base command when called without any subcommands
var (
	source, target, pattern string
	rootCmd                 = &cobra.Command{
		Use:   "filer",
		Short: "filer - Interactive file sorting REPL: keep or delete files using regex patterns",
		Long: `Usage: filer [-s SOURCE] [-t TARGET] [-p REGEX]

	Interactive file sorting REPL with regex filtering

	Options:
	-s, --source DIR    Source directory (default: current)
	-t, --target DIR    Target directory for kept files (default: keep in place)
	-p, --pattern REGEX Regular expression pattern to filter files

	Commands: [K]eep, [D]elete, [Q]uit
	Example: filer -s ~/Downloads -p "\.jpg$" -t ./images`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		Run: func(cmd *cobra.Command, args []string) {
			entries, err := os.ReadDir(source)
			if err != nil {
				fmt.Println(err)
				return
			}

			if pattern != "" {
				entries, err = filterAndSortFiles(entries)
				if err != nil {
					fmt.Println(err)
					return
				}
			}

			if target != "" {
				err := os.MkdirAll(target, 0755)
				if err != nil {
					fmt.Println(err)
					return
				}
			}

			oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
			if err != nil {
				fmt.Println("Error setting terminal to raw mode:", err)
				return
			}
			defer term.Restore(int(os.Stdin.Fd()), oldState)

			reader := bufio.NewReader(os.Stdin)

			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}

				fmt.Print("\033[2K\r")
				fmt.Printf("Name: %s [K]eep, [D]elete, [Q]uit?", entry.Name())

				char, _, err := reader.ReadRune()
				if err != nil {
					fmt.Println(err)
					return
				}

				switch strings.ToLower(string(char)) {
				case "q":
					return
				case "k":
					if target == "" {
						continue
					}
					err := moveFileSafe(source+"/"+entry.Name(), target+"/"+entry.Name())
					if err != nil {
						fmt.Println(err)
						return
					}
				case "d":
					err := os.Remove(source + "/" + entry.Name())
					if err != nil {
						fmt.Println(err)
						return
					}
				default:
					continue
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
	rootCmd.Flags().StringVarP(&pattern, "pattern", "p", "", "Regular expression pattern to filter files")
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

func filterAndSortFiles(entries []os.DirEntry) ([]os.DirEntry, error) {
	var filtered []os.DirEntry

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		matched, err := regexp.MatchString(pattern, entry.Name())
		if err != nil {
			return nil, err
		}
		if matched {
			filtered = append(filtered, entry)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Name() < filtered[j].Name()
	})

	return filtered, nil
}
