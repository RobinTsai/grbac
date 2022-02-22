package main

import (
	"grbac-gen/pkg/gen"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

const (
	searchDirFlag    = "dir"
	outputDirFlag    = "outDir"
	parseDepthFlag   = "parseDepth"
	excludeFilesFlag = "excludeFiles"
	outputFileFlag   = "output"
	formatFlag       = "format"
	tagFlag          = "tag"
	ssRoleFlag       = "ssRole"
	tidyFlag         = "tidy"
)

var initFlags = []cli.Flag{
	&cli.StringFlag{
		Name:  searchDirFlag,
		Value: "./",
		Usage: "Directory you want to parse",
	},
	&cli.StringFlag{
		Name:  outputDirFlag,
		Value: "./",
		Usage: "Directory you want to output.",
	},
	&cli.IntFlag{
		Name:  parseDepthFlag,
		Usage: "Depth you want to scan go files",
		Value: 3,
	},
	&cli.StringSliceFlag{
		Name:  excludeFilesFlag,
		Usage: "Which files/folders to be excludes",
		Value: nil,
	},
	&cli.StringFlag{
		Name:    outputFileFlag,
		Aliases: []string{"o"},
		Usage:   "What file name to be out",
		Value:   "permission",
	},
	&cli.StringFlag{
		Name:  formatFlag,
		Usage: "What format to be output. [*json/yaml]",
		Value: "json",
	},
	&cli.StringFlag{
		Name:     tagFlag,
		Usage:    "What tag",
		Required: true,
	},
	&cli.StringFlag{
		Name:  ssRoleFlag,
		Usage: "Set Super-Super-Role who has all permissions",
	},
	&cli.BoolFlag{
		Name:  tidyFlag,
		Usage: "tidy and compress rules",
	},
}

func main() {
	app := cli.NewApp() // 创建一个 cli 的 app
	app.Version = "0.0.1"
	app.Usage = "Automatically gen GRBAC configure json for Go project."
	app.Commands = []*cli.Command{ // app 注册子命令
		{
			Name:   "init",
			Usage:  "Create permission.json",
			Action: initAction, // 子命令的动作
			Flags:  initFlags,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

func initAction(c *cli.Context) error {
	return gen.New().Build(&gen.Config{
		SearchDir:    c.String(searchDirFlag),
		OutputDir:    c.String(outputDirFlag),
		ParseDepth:   c.Int(parseDepthFlag),
		ExcludeFiles: c.StringSlice(excludeFilesFlag),
		OutputFile:   c.String(outputFileFlag),
		Format:       c.String(formatFlag),
		Tag:          c.String(tagFlag),
		SsRole:       c.String(ssRoleFlag),
		Tidy:         c.Bool(tidyFlag),
	})
}
