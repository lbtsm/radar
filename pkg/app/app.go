package app

import (
	"fmt"
	"github.com/mapprotocol/filter/internal/pkg/constant"
	"github.com/urfave/cli/v2"
	"os"
)

type App struct {
	app cli.App
}

func New(cmds ...*cli.Command) *App {
	ret := App{}
	app := cli.NewApp()
	app.Copyright = "Copyright 2023 MAP Protocol 2023 Authors"
	app.Name = "filter"
	app.Usage = "Filter"
	app.Authors = []*cli.Author{{Name: "MAP Protocol 2023"}}
	app.Version = "1.0.0"
	app.EnableBashCompletion = true
	app.Flags = append(app.Flags, constant.ConfigFileFlag)
	app.Commands = append(app.Commands, cmds...)

	return &ret
}

func (a *App) Run() {
	if err := a.app.Run(os.Args); err != nil {
		fmt.Printf("%v \n", err)
		os.Exit(1)
	}
}
