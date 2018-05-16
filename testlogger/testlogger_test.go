package testlogger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bitrise-io/go-utils/log"
)

var dd = log.DefaultLogger{}

func TestTestFns(t *testing.T) {
	var b bytes.Buffer
	log.SetOutWriter(&b)

	model := MyPkg{}
	SetLogger(dd)
	model.TestFn1()
	model.TestFn2()

	res1 := b.String()
	var b2 bytes.Buffer
	log.SetOutWriter(&b2)

	log.TDonef("function called")
	log.Printf("some data")
	log.Errorf("+ some data")
	log.TWarnf("+ some data")
	log.Printf("some data")
	log.TDonef("function called")
	log.Errorf("+ some data")
	log.TWarnf("+ some data")

	res2 := b2.String()

	require.Equal(t, res1, res2)
}
