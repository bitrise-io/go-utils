package linepipe

import (
	"io"
	"strings"
	"testing"
)

// buildLog is a stand-in for what a build writes: colored progress, plain task lines and a long one.
var buildLog = []byte(strings.Repeat(
	"\x1b[1m> Task :app:compileDebugKotlin\x1b[0m\n"+
		"w: file:///src/app/Foo.kt:12:9: variable is never used\n"+
		"Download https://repo1.maven.org/maven2/org/example/lib/1.2.3/lib-1.2.3.jar\n"+
		strings.Repeat("x", 4000)+"\n",
	200))

func BenchmarkWriterForwardOnly(b *testing.B) {
	benchmarkWriter(b, New(io.Discard))
}

func BenchmarkWriterWithTapReadingText(b *testing.B) {
	tap := tapFunc(func(l *Line) { _ = l.Text() })
	benchmarkWriter(b, New(io.Discard, WithTap(tap)))
}

func BenchmarkWriterWithTapReadingRaw(b *testing.B) {
	tap := tapFunc(func(l *Line) { _ = len(l.Raw()) })
	benchmarkWriter(b, New(io.Discard, WithTap(tap)))
}

func benchmarkWriter(b *testing.B, w *Writer) {
	b.SetBytes(int64(len(buildLog)))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := w.Write(buildLog); err != nil {
			b.Fatal(err)
		}
	}
}
