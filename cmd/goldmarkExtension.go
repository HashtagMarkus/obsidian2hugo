package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type Ext struct {
	Title string
	Description string
	DescriptionTag string
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
				if h2Text == e.DescriptionTag {
					tldrFound = true
				}
			}
		}
		// Extract description
		if tldrFound && c.Kind() == gast.KindParagraph {
			log.Debug(c.FirstChild().Kind().String())
			e.Description = e.dumpStr(c, reader.Source(), "")
			return
		}
	}
}

func (e *Ext) dumpStr(c gast.Node, source []byte, str string) string {
	res := str
	for l := c.FirstChild(); l != nil; l = l.NextSibling() {
		if l.Kind() == gast.KindText {
			desc := string(l.Text(source))
			log.Debug(desc)
			res = res + desc
		} else {
			if l.HasChildren() {
				res = e.dumpStr(l, source, res)
			}
		}
	}
	return res
}
