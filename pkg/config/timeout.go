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
	RuntimeRequestTimeout = 1 * Minute
	RPCRequestTimeout     = 1 * Minute
)

const (
	// IdenticalErrorDelay How frequently to report identical errors
	IdenticalErrorDelay = 1 * time.Minute

	MaxBackoffDelay      = 3 * time.Second
	BaseBackoffDelay     = 100 * time.Millisecond
	MinConnectionTimeout = 5 * time.Second
)
