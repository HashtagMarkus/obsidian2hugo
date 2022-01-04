module github.com/HashtagMarkus/obsidian2hugo

go 1.17

require (
	github.com/adrg/frontmatter v0.2.0
	github.com/spf13/cobra v1.3.0
	github.com/HashtagMarkus/obsidian2hugo/cmd v0.0.0
)

replace "github.com/HashtagMarkus/obsidian2hugo/cmd" v0.0.0 => "/home/zerd/go/src/github.com/HashtagMarkus/obsidian2hugo/cmd"

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
