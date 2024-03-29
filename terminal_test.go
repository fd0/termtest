package termtest_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/fd0/termtest"
)

var termTests = []struct {
	x, y    int
	command string
	output  []string
}{
	{
		x:       10,
		y:       10,
		command: "echo AAAAAAAAAA",
		output: []string{
			"AAAAAAAAAA",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		},
	},
	{
		x:       10,
		y:       10,
		command: "echo AAAAAAAAAAB",
		output: []string{
			"AAAAAAAAAA",
			"B",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
			"",
		},
	},
}

func TestTerminal(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	term, err := termtest.New(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	for i, test := range termTests {
		buf, err := term.Run(ctx, test.x, test.y, test.command)
		if err != nil {
			t.Errorf("test %d failed, error %v", i, err)
			continue
		}

		want := ""
		for _, s := range test.output {
			want += s + "\n"
		}

		if string(buf) != want {
			t.Errorf("test %d failed, want output:\n  %q\ngot:\n  %q", i,
				want, string(buf))
		}
	}

	if err = term.Exit(ctx); err != nil {
		t.Fatal(err)
	}
}

func ExampleTerminal() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create terminal
	term, err := termtest.New(ctx, nil)
	if err != nil {
		panic(err)
	}

	// run something on a terminal which is 12 characters wide and has 13 lines
	buf, err := term.Run(ctx, 12, 13, "echo This is a long line that will wrap for sure")
	if err != nil {
		panic(err)
	}

	// print the output string
	fmt.Printf("%q\n", buf)

	// stop the tmux instance
	err = term.Exit(ctx)
	if err != nil {
		panic(err)
	}

	// Output: "This is a lo\nng line that\n will wrap f\nor sure\n\n\n\n\n\n\n\n\n\n"
}
