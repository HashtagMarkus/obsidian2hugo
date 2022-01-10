/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"bytes"
	parser2 "github.com/HashtagMarkus/obsidian2hugo/cmd/parser"
	"github.com/HashtagMarkus/obsidian2hugo/cmd/parser/pageparser"
	"github.com/adrg/frontmatter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type FrontMatter struct {
	Published bool `yaml:"published"`
	Date string `yaml:"date"`
	Tags []string `yaml:"tags"`
}

var rootCmd = &cobra.Command{
	Use:   "obsidian2hugo",
	Short: "Finds obsidian pages marked as `published: true` and copies the files into the hugo directory",
	Long: `This tool was created to be able to export blog posts created inside obsidian for the usage inside a hugo blog.`,

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
		descriptionTag, _ := cmd.Flags().GetString("descriptionSection")
		keepTitle, _ := cmd.Flags().GetBool("keepTitle")

		files, err := WalkMatch(source, "*.md")
		for _, file := range files {
			// Read frontmatter stuff
			f, err := os.Open(file)
			defer f.Close()
			if err != nil {
				log.Fatal(err)
			}
			buf := bufio.NewReader(f)
			matter := FrontMatter{}
			frontmatter.Parse(buf, &matter)
			if matter.Published {
				relativeFileFromSource := strings.TrimPrefix(file, source)
				folder := filepath.Dir(relativeFileFromSource)

				// tl;dr section for Frontmatter description and h1 tag for frontatter title
				e := Ext{DescriptionTag: descriptionTag}
				gm := goldmark.New(goldmark.WithExtensions(&e))
				bufarr, _ := ioutil.ReadFile(file)

				reader := text.NewReader(bufarr)
				gm.Parser().Parse(reader)

				// Reset file to parse frontmatter section
				f.Seek(0,0)
				res, err := pageparser.ParseFrontMatterAndContent(f)
				if err != nil {
					log.Error(err)
				}
				if len(e.Description) > 0 { // If found, set description in frontmatter section
					res.FrontMatter["description"] = e.Description
				}
				if len(e.Title) > 0 { // If h1 header found, set in frontmatter section
					res.FrontMatter["title"] = strings.TrimRight(e.Title, "\n")
					if !keepTitle {
						res.Content = []byte(strings.Replace(string(res.Content), "# " + e.Title, "", 1))
					}
				}

				// Copy content to target
				CopyDir(path.Join(source, folder), path.Join(target, folder))

				// Save file with adjusted frontmatter tags and removed h1 title in target dir
				f2, err := os.OpenFile(path.Join(target, folder, filepath.Base(f.Name())), os.O_RDWR, 0755)
				defer f2.Close()
				log.Debug(res.FrontMatter)

				var writeBuf bytes.Buffer
				if len(res.FrontMatter) != 0 {
					err := parser2.InterfaceToFrontMatter(res.FrontMatter, "toml", &writeBuf)
					if err != nil {
						log.Error(err)
					}
				}

				writeBuf.WriteString(string(res.Content))
				f2.Truncate(0)
				f2.Seek(0,0)
				f2.Write(writeBuf.Bytes())
			}
		}
	 },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("source", "s", ".", "Source to obsidian markdown files (root of blog posts tree, e.g.: <obsidianvault>/blogposts)")
	rootCmd.Flags().StringP("destination", "d", "", "Destination of hugo posts folder (e.g. <hugoroot>/content/posts)")
	rootCmd.Flags().BoolP("keepTitle", "k", false, "Don't delete h1 header after frontmatter extraction")
	rootCmd.Flags().StringP("descriptionSection", "t", "tl;dr", "The content below this h2 header is used as the frontmatter description")

	rootCmd.MarkFlagRequired("destination")
}


