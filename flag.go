package main

import "flag"

const (
	FlagNameTemplate = "template"
)

type FlagMap map[string]*string

func mustParseFlags() FlagMap {
	availableFlags := FlagMap{
		FlagNameTemplate: flag.String(FlagNameTemplate, "random", "Set a template"),
	}
	flag.Parse()
	return availableFlags
}
