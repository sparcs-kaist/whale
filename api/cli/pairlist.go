package cli

import (
	"github.com/sparcs-kaist/whale"

	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"strings"
)

type pairList []whale.Pair

// Set implementation for a list of whale.Pair
func (l *pairList) Set(value string) error {
	parts := strings.SplitN(value, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("expected NAME=VALUE got '%s'", value)
	}
	p := new(whale.Pair)
	p.Name = parts[0]
	p.Value = parts[1]
	*l = append(*l, *p)
	return nil
}

// String implementation for a list of pair
func (l *pairList) String() string {
	return ""
}

// IsCumulative implementation for a list of pair
func (l *pairList) IsCumulative() bool {
	return true
}

func pairs(s kingpin.Settings) (target *[]whale.Pair) {
	target = new([]whale.Pair)
	s.SetValue((*pairList)(target))
	return
}
