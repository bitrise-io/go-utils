#!/bin/bash
# Exits cleanly when signalled, so the test can tell a handled SIGTERM from a hard kill.
trap 'exit 0' TERM
while true; do sleep 0.1; done
