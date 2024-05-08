package main

import (
	"github.com/mapprotocol/filter/internal/api"
	"github.com/mapprotocol/filter/internal/filter"
	"github.com/mapprotocol/filter/pkg/app"
)

func main() {
	a := app.New(filter.Command, api.Command)
	a.Run()
}
