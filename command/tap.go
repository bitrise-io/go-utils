package command

import "io"

// observer is the second route of a command's output. It is shown every byte the caller's writer
// took, after the writer took it, so it can read the output but never hold it back or change it.
type observer interface {
	observe(p []byte)
}

// tap sits between the tool and the caller's writer. It is the one place output is forwarded.
type tap struct {
	dest     io.Writer
	observer observer
}

// Write forwards p untouched, then shows the observer what was written.
func (t *tap) Write(p []byte) (int, error) {
	n, err := t.dest.Write(p)
	// A count outside p is passed back unused, so that io.Copy reports it as an invalid write.
	if n < 0 || n > len(p) {
		return n, err
	}

	if n > 0 {
		t.observer.observe(p[:n])
	}

	return n, err
}
