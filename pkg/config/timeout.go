package config

import (
	"time"
)

const (
	minDuration int64 = -1 << 63
	maxDuration int64 = 1<<63 - 1
)

const (
	Nanosecond  int64 = 1
	Microsecond       = 1000 * Nanosecond
	Millisecond       = 1000 * Microsecond
	Second            = 1000 * Millisecond
	Minute            = 60 * Second
	Hour              = 60 * Minute
)

const (
	RuntimeRequestTimeout = 30 * Minute
)

const (
	// How frequently to report identical errors
	IdenticalErrorDelay = 1 * time.Minute

	// connection parameters
	MaxBackoffDelay      = 3 * time.Second
	BaseBackoffDelay     = 100 * time.Millisecond
	MinConnectionTimeout = 5 * time.Second
)
