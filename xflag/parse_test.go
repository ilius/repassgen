package xflag

import (
	"flag"
	"io"
	"reflect"
	"testing"
)

func TestParseToEnd(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		fs := flag.NewFlagSet("name", flag.ContinueOnError)
		debugP := fs.Bool("debug", false, "debug mode")
		requireNoError(t, ParseToEnd(fs, []string{}))
		requireFalse(t, *debugP)
		requireEqual(t, 0, fs.NFlag())
	})
	t.Run("no args", func(t *testing.T) {
		fs := flag.NewFlagSet("name", flag.ContinueOnError)
		debugP := fs.Bool("debug", false, "debug mode")
		err := ParseToEnd(fs, []string{"--debug", "true"})
		requireNoError(t, err)
		requireEqual(t, 1, fs.NFlag())
		requireTrue(t, *debugP)
	})
	t.Run("before with args", func(t *testing.T) {
		fs := flag.NewFlagSet("name", flag.ContinueOnError)
		debugP := fs.Bool("debug", false, "debug mode")
		err := ParseToEnd(fs, []string{"--debug=true", "arg1", "arg2"})
		requireNoError(t, err)
		requireTrue(t, *debugP)
		requireEqual(t, 1, fs.NFlag())
		requireEqual(t, []string{"arg1", "arg2"}, fs.Args())
	})
	t.Run("after with args", func(t *testing.T) {
		fs := flag.NewFlagSet("name", flag.ContinueOnError)
		debugP := fs.Bool("debug", false, "debug mode")
		err := ParseToEnd(fs, []string{"arg1", "arg2", "--debug"})
		requireNoError(t, err)
		requireTrue(t, *debugP)

		f := fs.Lookup("debug")
		requireEqual(t, "true", f.Value.String())
		requireEqual(t, 1, fs.NFlag())
		requireEqual(t, []string{"arg1", "arg2"}, fs.Args())
	})
	t.Run("before and after with args", func(t *testing.T) {
		fs, c := newFlagset()
		args := []string{
			"--flag1=value1",
			"--flag3=true",
			"arg1",
			"arg2",
			"--flag2=value2",
			"--flag4=false",
			"arg3",
		}
		err := ParseToEnd(fs, args)
		requireNoError(t, err)
		requireEqual(t, config{
			flag1: "value1",
			flag3: true,
			flag2: "value2",
			flag4: false,
		}, *c)
		requireEqual(t, []string{"arg1", "arg2", "arg3"}, fs.Args())
		requireEqual(t, 4, fs.NFlag())
	})
	t.Run("break", func(t *testing.T) {
		fs, c := newFlagset()
		args := []string{
			"--flag1=value1",
			"--flag3=true",
			"arg1",
			"arg2",
			"--flag2=value2",
			"--flag4=false",
			"arg3",
			"--",
			"arg4",
			"arg5",
			"--flag4=true", // This is now a positional argument no matter what.
		}
		err := ParseToEnd(fs, args)
		requireNoError(t, err)
		requireEqual(t, config{
			flag1: "value1",
			flag3: true,
			flag2: "value2",
			flag4: false,
		}, *c)
		requireEqual(t, []string{"arg1", "arg2", "arg3", "--", "arg4", "arg5", "--flag4=true"}, fs.Args())
		requireEqual(t, 4, fs.NFlag())
	})
	t.Run("unknown flag before", func(t *testing.T) {
		fs, _ := newFlagset()
		args := []string{
			"--flag1=value1",
			"--some-unknown-flag=foo", // This gets treated as a flag.
			"arg1",
		}
		err := ParseToEnd(fs, args)
		requireError(t, err)
		requireEqual(t, err.Error(), "flag provided but not defined: -some-unknown-flag")
	})
	t.Run("unknown flag after", func(t *testing.T) {
		fs, _ := newFlagset()
		args := []string{
			"--flag1=value1",
			"arg1",
			"--some-unknown-flag=foo", // This gets treated as a flag.
		}
		err := ParseToEnd(fs, args)
		requireError(t, err)
		requireEqual(t, err.Error(), "flag provided but not defined: -some-unknown-flag")
	})
}

type config struct {
	flag1 string
	flag2 string
	flag3 bool
	flag4 bool
}

func newFlagset() (*flag.FlagSet, *config) {
	fs := flag.NewFlagSet("name", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	c := new(config)
	fs.StringVar(&c.flag1, "flag1", "asdf", "flag1 usage")
	fs.StringVar(&c.flag2, "flag2", "qwerty", "flag2 usage")
	fs.BoolVar(&c.flag3, "flag3", false, "flag3 usage")
	fs.BoolVar(&c.flag4, "flag4", true, "flag4 urage")
	return fs, c
}

func requireError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error")
	}
}

func requireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func requireTrue(t *testing.T, b bool) {
	t.Helper()
	if !b {
		t.Fatal("expected true")
	}
}

func requireFalse(t *testing.T, b bool) {
	t.Helper()
	if b {
		t.Fatal("expected false")
	}
}

func requireEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}
