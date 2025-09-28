package main

import "fmt"

type chartStyle int

const (
	csBarSimple chartStyle = iota
	csBarHorizontal1px
	csBarVertical1px
)

func (cs chartStyle) String() string {
	switch cs {
	case csBarSimple:
		return "bar-simple"
	case csBarHorizontal1px:
		return "bar-horizontal-1px"
	case csBarVertical1px:
		return "bar-vertical-1px"
	default:
		return "unknown"
	}
}

func (cs *chartStyle) Set(s string) error {
	switch s {
	case "bar-simple":
		*cs = csBarSimple
	case "bar-horizontal-1px":
		*cs = csBarHorizontal1px
	case "bar-vertical-1px":
		*cs = csBarVertical1px
	default:
		return fmt.Errorf("unknown chart style %q", s)
	}
	return nil
}
