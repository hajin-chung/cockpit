package main

import (
	"time"

	"github.com/oklog/ulid/v2"
)

const MAX_ID string = "7ZZZZZZZZZZZZZZZZZZZZZZZZZ"

func IdGen() string {
	return ulid.Make().String()
}

func FormatNow() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}
