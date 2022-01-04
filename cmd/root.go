/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"fmt"
	"github.com/adrg/frontmatter"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func CopyFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

// Dir copies a whole directory recursively
func CopyDir(src string, dst string) error {

	log.Printf("Copy %s to %s\n\n", src, dst)

	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = CopyDir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = CopyFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func WalkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

type FrontMatter struct {
	Published bool `yaml:"published"`
	Date string `yaml:"date"`
	Tags []string `yaml:"tags"`
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "obsidian2hugo",
	Short: "Finds obsidian pages marked as `published: true` and copies the files into the hugo directory",
	Long: `TODO`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	 Run: func(cmd *cobra.Command, args []string) {
		source, err := cmd.Flags().GetString("source")
		if err != nil {
			log.Fatal("`source` flag missing")
		}
		target, err := cmd.Flags().GetString("destination")
		if err != nil {
		 	log.Fatal("`target` flag missing")
		}
		source, err = filepath.Abs(source)
		if err != nil {
			log.Fatal("Cannot get absolute file path from `source`")
		}
		files, err := WalkMatch(source, "*.md")
		for _, file := range files {
			// Read frontmatter stuff
			f, err := os.Open(file)
			if err != nil {
				log.Fatal(err)
			}
			buf := bufio.NewReader(f)
			matter := FrontMatter{}
			frontmatter.Parse(buf, &matter)
			if matter.Published {
				relativeFileFromSource := strings.TrimPrefix(file, source)
				log.Println(relativeFileFromSource)
				folder := filepath.Dir(relativeFileFromSource)
				log.Println(folder)

				// TODO: copy all files inside this folder to target tree
				CopyDir(path.Join(source, folder), path.Join(target, folder))

				// TODO: do I need to patch anything else?
			}
		}
	 },
}

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

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.obsidian2hugo.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringP("source", "s", ".", "Source to obsidian markdown files")
	rootCmd.Flags().StringP("destination", "d", "", "Destination of hugo folder")

	rootCmd.MarkFlagRequired("destination")
}


