package math

import (
	builtinMath "math"

	logger "github.com/multiversx/mx-chain-logger-go"
)

var log = logger.GetOrCreate("vm/overflow")

// AddUint64 performs addition on uint64 and logs an error if the addition overflows.
//
// ISSUE-079 (partial — log-level step only; full sentinel-pattern audit
// across callers is tracked separately): promoted log level from Trace
// to Warn so overflow events surface in default-level production logs.
// On overflow returns MaxUint64 sentinel. Callers that need to
// distinguish "valid MaxUint64" from "overflow occurred" must use
// AddUint64WithErr instead.
func AddUint64(a, b uint64) uint64 {
	res, err := AddUint64WithErr(a, b)
	if err != nil {
		log.Warn("AddUint64 overflow — saturating to MaxUint64", "a", a, "b", b)
		return builtinMath.MaxUint64
	}

	return res
}

// AddUint64WithErr performs addition on uint64 and returns an error if the addition overflows
func AddUint64WithErr(a, b uint64) (uint64, error) {
	s := a + b
	if s >= a && s >= b {
		return s, nil
	}

	return builtinMath.MaxUint64, ErrAdditionOverflow
}

// SubUint64 performs subtraction on uint64, in case of underflow returns 0
func SubUint64(a, b uint64) uint64 {
	if a < b {
		return 0
	}

	return a - b
}

// MulUint64 performs multiplication on uint64 and logs an error if the multiplication overflows.
//
// ISSUE-079 (partial — log-level step only; see AddUint64 doc-comment).
func MulUint64(a, b uint64) uint64 {
	res, err := MulUint64WithErr(a, b)
	if err != nil {
		log.Warn("MulUint64 overflow — saturating to MaxUint64", "a", a, "b", b)
		return builtinMath.MaxUint64
	}

	return res
}

// MulUint64WithErr performs multiplication on uint64 and returns an error if the multiplication overflows
func MulUint64WithErr(a, b uint64) (uint64, error) {
	res := a * b
	if a == 0 || b == 0 || a == res/b {
		return res, nil
	}

	return builtinMath.MaxUint64, ErrMultiplicationOverflow
}

// AddInt32 performs addition on int32 and logs an error if the addition overflows
func AddInt32(a, b int32) int32 {
	res, err := AddInt32WithErr(a, b)
	if err != nil {
		log.Trace("AddInt32 overflow", "a", a, "b", b)
		return builtinMath.MaxInt32
	}

	return res
}

// AddInt32WithErr performs addition on int32 and returns an error if the addition overflows
func AddInt32WithErr(a, b int32) (int32, error) {
	s := a + b
	if (s > a) == (b > 0) {
		return s, nil
	}

	return builtinMath.MaxInt32, ErrAdditionOverflow
}

// SubInt performs subtraction on int and logs an error if the subtraction overflows.
//
// ISSUE-086: on underflow, returns the MinInt64 sentinel. Promoted log
// level from Trace to Warn so overflow events are visible in default-
// level production logs (Trace is typically disabled in production).
// Callers that need to distinguish "valid MinInt64" from "underflow
// occurred" must check inputs themselves before calling this function.
func SubInt(a, b int) int {
	res := a - b
	if (res < a) == (b > 0) {
		return res
	}

	log.Warn("SubInt underflow — saturating to MinInt64", "a", a, "b", b)
	return builtinMath.MinInt64
}
