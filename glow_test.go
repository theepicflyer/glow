package main

import (
	"testing"
)

func TestGlowFlags(t *testing.T) {
	tt := []struct {
		args  []string
		check func() bool
	}{
		{
			args: []string{"-p"},
			check: func() bool {
				return pager
			},
		},
		{
			args: []string{"-s", "light"},
			check: func() bool {
				return style == "light"
			},
		},
		{
			args: []string{"-w", "40"},
			check: func() bool {
				return width == 40
			},
		},
	}

	for _, v := range tt {
		err := rootCmd.ParseFlags(v.args)
		if err != nil {
			t.Fatal(err)
		}
		if !v.check() {
			t.Errorf("Parsing flag failed: %s", v.args)
		}
	}
}

func TestResolveWidth(t *testing.T) {
	tt := []struct {
		name          string
		configured    uint
		isTerminal    bool
		terminalWidth uint
		want          uint
	}{
		{"adaptive follows terminal width", 0, true, 200, 200},
		{"adaptive is uncapped on wide terminals", 0, true, 300, 300},
		{"adaptive falls back when not a terminal", 0, false, 0, defaultWidth},
		{"adaptive falls back when terminal size is unknown", 0, true, 0, defaultWidth},
		{"explicit configured width is respected", 100, true, 200, 100},
		{"explicit configured width is respected without a terminal", 100, false, 0, 100},
	}

	for _, v := range tt {
		t.Run(v.name, func(t *testing.T) {
			if got := resolveWidth(v.configured, v.isTerminal, v.terminalWidth); got != v.want {
				t.Errorf("resolveWidth(%d, %v, %d) = %d, want %d", v.configured, v.isTerminal, v.terminalWidth, got, v.want)
			}
		})
	}
}
