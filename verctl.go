//
// Copyright (c) 2016 ЗАО Геликон Про http://www.gelicon.biz
//
package main

import (
	"os"
	"fmt"
	"flag"
	"bufio"
	"io/ioutil"
	"strings"
	"strconv"
	"./version"
)

var (
	Version = "unknown"
)

type Params struct {
	UseTool   bool
	UseGen    bool
	UseVer    bool
	UseSub    bool
	UseLast   bool
	Text      string
	FileName  string
	ShowLevel bool
	Release   bool
	Trunk     bool
}

func BindParams(flags *flag.FlagSet) *Params {

	r := Params{}

	flags.BoolVar(&r.UseTool, "tool", false, "Print tool version")

	flags.BoolVar(&r.ShowLevel, "level", false, "Print version level")

	flags.BoolVar(&r.UseGen, "gen", false, "Use generation level")
	flags.BoolVar(&r.UseVer, "ver", false, "Use version level")
	flags.BoolVar(&r.UseSub, "sub", false, "Use subversion level")
	flags.BoolVar(&r.UseLast, "last", false, "Use last level")

	flags.BoolVar(&r.Release, "release", false, "Make release")
	flags.BoolVar(&r.Trunk, "trunk", false, "Trunk version string")

	flags.StringVar(&r.Text, "text", "<stdin>", "Use text from command line instead of standard stream")
	flags.StringVar(&r.FileName, "file", "", "Use a file instead of standard stream")

	return &r
}

func RunTool(params *Params) int {

	ver := version.Ver{}

	switch {
	case params.UseGen:
		ver.Start(1)

	case params.UseVer:
		ver.Start(2)

	case params.UseSub:
		ver.Start(3)

	default:
		ver.Start(1)
	}

	var text string

	switch {
	case params.Text != "<stdin>":
		text = params.Text

	case params.UseTool:
		text = Version

	case params.FileName == "":
		scanner := bufio.NewScanner(os.Stdin)

		if !scanner.Scan() {
			return 1
		}

		text = scanner.Text()

	default:
		data, _ := ioutil.ReadFile(params.FileName)
		text = string(data)
	}

	text = strings.TrimSpace(text)

	if text != "" {
		ver.Parse(text)

		switch {
		case params.Trunk || params.UseTool:
			switch {
			case params.UseGen:
				ver.Trunk(1)

			case params.UseVer:
				ver.Trunk(2)

			case params.UseSub:
				ver.Trunk(3)

			case params.UseLast:
				ver.Trunk(ver.Level())
			}
		default:
			switch {
			case params.UseGen:
				ver.Next(1)

			case params.UseVer:
				ver.Next(2)

			case params.UseSub:
				ver.Next(3)

			case params.UseLast:
				ver.Next(ver.Level())
			}
		}
	}

	if params.Release {
		ver.Release()
	}

	var output string

	switch {
	case params.ShowLevel:
		output = strconv.Itoa(ver.Level())

	default:
		output = ver.Render()
	}

	switch {
	case params.FileName == "":
		fmt.Println(output)

	default:
		if err := ioutil.WriteFile(params.FileName, []byte(output), 0); err != nil {
			return 1
		}
	}

	return 0
}

func main() {

	params := BindParams(flag.CommandLine)

	flag.Parse()

	os.Exit(RunTool(params))
}
