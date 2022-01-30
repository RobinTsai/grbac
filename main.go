package main

import (
	"grbac-gen/pkg/gen"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

const (
	searchDirFlag  = "dir"
	outputDirFlag  = "output"
	parseDepthFlag = "parseDepth"
	mainFileFlag   = "mainFile"
	outputVarFlag  = "variable"
)

var initFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    searchDirFlag,
		Aliases: []string{"d"},
		Value:   "./",
		Usage:   "Directory you want to parse",
	},
	&cli.StringFlag{
		Name:    outputDirFlag,
		Aliases: []string{"o"},
		Value:   "./docs",
		Usage:   "Directory you want to output.",
	},
	&cli.IntFlag{
		Name:  parseDepthFlag,
		Usage: "Depth you want to scan go files",
		Value: 3,
	},
	&cli.StringFlag{
		Name:  mainFileFlag,
		Value: "main.go",
		Usage: "Main file name",
	},
	&cli.StringFlag{
		Name:    outputVarFlag,
		Aliases: []string{"v"},
		Value:   "permissionStr",
		Usage:   "The variable of permission config string",
	},
}

func main() {
	app := cli.NewApp() // 创建一个 cli 的 app
	app.Version = "0.0.1"
	app.Usage = "Automatically gen GRBAC configure json for Go project."
	app.Commands = []*cli.Command{ // app 注册子命令
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "Create permission.json",
			Action:  initAction, // 子命令的动作
			Flags:   initFlags,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

func initAction(c *cli.Context) error {
	return gen.New().Build(&gen.Config{
		SearchDir:  c.String(searchDirFlag),
		OutputDir:  c.String(outputDirFlag),
		ParseDepth: c.Int(parseDepthFlag),
		MainFile:   c.String(mainFileFlag),
		OutputVar:  c.String(outputVarFlag),
	})
}
