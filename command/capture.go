package command

import "fmt"

// stderrCaptureLimit is how much of each end of stderr is kept, the same as os/exec keeps for
// Cmd.Output.
const stderrCaptureLimit = 32 << 10

// headTailBuffer keeps the first and the last limit bytes written to it and drops the middle, the
// way os/exec captures stderr for Cmd.Output. A command that prints megabytes to stderr must not be
// held in memory for the sake of an error message.
type headTailBuffer struct {
	limit   int
	head    []byte
	tail    []byte // ring buffer once it is full
	tailOff int
	skipped int64
}

func (b *headTailBuffer) Write(p []byte) (int, error) {
	n := len(p)

	if room := b.limit - len(b.head); room > 0 {
		take := min(room, len(p))
		b.head = append(b.head, p[:take]...)
		p = p[take:]
	}

	// Only the last limit bytes of p can end up in the tail.
	if over := len(p) - b.limit; over > 0 {
		b.skipped += int64(over)
		p = p[over:]
	}

	if room := b.limit - len(b.tail); room > 0 {
		take := min(room, len(p))
		b.tail = append(b.tail, p[:take]...)
		p = p[take:]
	}

	for len(p) > 0 {
		copied := copy(b.tail[b.tailOff:], p)
		p = p[copied:]
		b.skipped += int64(copied)
		b.tailOff = (b.tailOff + copied) % b.limit
	}

	return n, nil
}

// Bytes returns what was kept, with a note in the middle saying how much was dropped.
func (b *headTailBuffer) Bytes() []byte {
	if b.skipped == 0 {
		return append(append([]byte(nil), b.head...), b.tail...)
	}

	out := append([]byte(nil), b.head...)
	out = append(out, fmt.Sprintf("\n... omitting %d bytes ...\n", b.skipped)...)
	out = append(out, b.tail[b.tailOff:]...)
	out = append(out, b.tail[:b.tailOff]...)

	return out
}
