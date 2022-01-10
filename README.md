# obsidian2hugo

## Usage

```bash
obsidian2hugo --help
This tool was created to be able to export blog posts created inside obsidian for the usage inside a hugo blog.

Usage:
  obsidian2hugo [flags]

Flags:
  -t, --descriptionSection string   The content below this h2 header is used as the frontmatter description (default "tl;dr")
  -d, --destination string          Destination of hugo posts folder (e.g. <hugoroot>/content/posts)
  -h, --help                        help for obsidian2hugo
  -k, --keepTitle                   Don't delete h1 header after frontmatter extraction
  -s, --source string               Source to obsidian markdown files (root of blog posts tree, e.g.: <obsidianvault>/blogposts) (default ".")
```

## Description

Read more about this tool in my [blog](https://task2.net/posts/2022-01-10-obsidian2hugo-exporter/2022-01-10-obsidian2hugo-exporter/)

## Additional notes

The code under the subdirectory `cmd/parser` was taken from the [hugo repository](https://github.com/gohugoio/hugo) in compliance with the Apache License 2.0. The code under said subdirectory was not changed by the author of obsidian2hugo.