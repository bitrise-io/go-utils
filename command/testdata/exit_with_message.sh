#!/usr/bin/env bash
# These messages are written to stdout
echo Error: first error
echo Error: second error
echo This is not an stdout error
# These messages are written to stderr
echo Error: third error >&2
echo Error: fourth error >&2
exit 1