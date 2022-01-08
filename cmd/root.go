/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	parser2 "github.com/HashtagMarkus/obsidian2hugo/cmd/parser"
	"github.com/HashtagMarkus/obsidian2hugo/cmd/parser/pageparser"
	"github.com/adrg/frontmatter"
	"github.com/spf13/cobra"
	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"io/ioutil"
	"log"
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

type Ext struct {
	Title string
	Description string
}

func (e *Ext) Extend(md goldmark.Markdown) {
	md.Parser().AddOptions(
		parser.WithASTTransformers(util.Prioritized(e, 999)),
	)
}

func (e* Ext) Transform(node *gast.Document, reader text.Reader, pc parser.Context) {
	tldrFound := false
	//node.Dump(reader.Source(), 0)
	for c :=  node.FirstChild(); c != nil; c = c.NextSibling() {
		if c.Kind() == gast.KindHeading {
			h := c.(*gast.Heading)
			if h.Level == 1 {
				e.Title = string(c.Text(reader.Source()))
			} else if h.Level == 2 {
				t := c.FirstChild().(*gast.Text)
				h2Text := string(t.Text(reader.Source()))
				if h2Text == "tl;dr" {
					tldrFound = true
				}
			}
		}
		// Extract description
		if tldrFound && c.Kind() == gast.KindParagraph {
			fmt.Println(c.FirstChild().Kind().String())
			e.Description = DumpStr(c, reader.Source(), "")
			return
		}
	}
}

func DumpStr(c gast.Node, source []byte, str string) string {
	res := str
	for l := c.FirstChild(); l != nil; l = l.NextSibling() {
		if l.Kind() == gast.KindText {
			desc := string(l.Text(source))
			fmt.Println(desc)
			res = res + desc
		} else {
			if l.HasChildren() {
				res = DumpStr(l, source, res)
			}
		}
	}
	return res
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
			defer f.Close()
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

				// Extract Frontmatter description
				e := Ext{}
				gm := goldmark.New(goldmark.WithExtensions(&e))
				bufarr, _ := ioutil.ReadFile(file)

				reader := text.NewReader(bufarr)
				gm.Parser().Parse(reader)

				// Add frontmatter stuff
				f.Seek(0,0)
				res, err := pageparser.ParseFrontMatterAndContent(f)
				if err != nil {
					fmt.Println(err)
				}
				if len(e.Description) > 0 {
					res.FrontMatter["description"] = e.Description // Add Desciption
				}
				if len(e.Title) > 0 {
					res.FrontMatter["title"] = strings.TrimRight(e.Title, "\n")
					res.Content = []byte(strings.Replace(string(res.Content), "# " + e.Title, "", 1))
				}

				// Copy content to target
				CopyDir(path.Join(source, folder), path.Join(target, folder))


				f2, err := os.OpenFile(path.Join(target, folder, filepath.Base(f.Name())), os.O_RDWR, 0755)
				fmt.Println(res.FrontMatter)

				var writeBuf bytes.Buffer
				if len(res.FrontMatter) != 0 {
					err := parser2.InterfaceToFrontMatter(res.FrontMatter, "toml", &writeBuf)
					if err != nil {

					}
				}

				writeBuf.WriteString(string(res.Content))
				//fmt.Println(writeBuf.String())
				f.Close()
				f2.Truncate(0)
				f2.Seek(0,0)
				f2.Write(writeBuf.Bytes())
				f2.Close()

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
	rootCmd.Flags().StringP("source", "s", ".", "Source to obsidian markdown files")
	rootCmd.Flags().StringP("destination", "d", "", "Destination of hugo folder")

	rootCmd.MarkFlagRequired("destination")
}


