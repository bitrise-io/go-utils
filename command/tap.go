package command

import "io"

// OutputObserver is shown the output of a command after it was forwarded to the caller's writer, so
// it can read the output but never hold it back or change it. When the command fails, what the
// observer reports becomes the reason of the error.
//
// An observer belongs to one output stream of one command. See Opts.NewOutputObserver for which
// stream that is.
type OutputObserver interface {
	// Observe is given what the command wrote, in the chunks it arrived in. A chunk can end in the
	// middle of a line.
	Observe(p []byte)
	// Close is called once, after the command has ended and everything it wrote was observed.
	Close()
	// Reasons returns the lines that say why the command failed, most telling first, or nothing
	// when the output held none. It is called after Close.
	Reasons() []string
}

// tap sits between the tool and the caller's writer. It is the one place output is forwarded.
type tap struct {
	dest     io.Writer
	observer OutputObserver
}

// Write forwards p untouched, then shows the observer what was written.
func (t *tap) Write(p []byte) (int, error) {
	n, err := t.dest.Write(p)
	// A count outside p is passed back unused, so that io.Copy reports it as an invalid write.
	if n < 0 || n > len(p) {
		return n, err
	}

	if n > 0 {
		t.observer.Observe(p[:n])
	}

	return n, err
}
