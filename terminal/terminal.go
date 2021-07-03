package terminal

import (
	"context"
	"os"
	"sync"

	"github.com/pkg/errors"

	"golang.org/x/term"
)

const prompt = "> "

type Terminal struct {
	t         *term.Terminal
	prevState *term.State
	prevFd    int
}

func Init() (*Terminal, error) {
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		return nil, errors.New("Not a terminal.")
	}

	prevState, err := term.MakeRaw(fd)
	if err != nil {
		return nil, errors.Errorf("Error setting the terminal to raw: %v", err)
	}

	return &Terminal{
		t:         term.NewTerminal(os.Stdin, prompt),
		prevState: prevState,
		prevFd:    fd,
	}, nil
}

func (t *Terminal) Run(ctx context.Context, wg sync.WaitGroup, print chan string) {
	var (
		msg, line string
		input     = make(chan string)
	)
	go t.reader(ctx, input)

	for {
		select {
		case msg = <-print:
			t.print(msg)

		case line = <-input:
			t.print(line)

		case <-ctx.Done():
			t.print("Exiting...")
			term.Restore(int(os.Stderr.Fd()), t.prevState)
			wg.Done()
			return
		}
	}
}

func (t *Terminal) reader(ctx context.Context, input chan string) {
	for {
		select {
		case <-ctx.Done():
			return

		default:
			line, _ := t.t.ReadLine()
			input <- line
		}
	}
}

func (t *Terminal) print(input string) {
	t.t.Write([]byte(input + "\n"))
}
