package internal

import (
	"encoding/json"
	"math/rand"
)

// D safely dereferences a value
func D[T comparable](v *T) T {
	var zero T
	if v == nil {
		return zero
	}
	return *v
}

// R returns a pointer value to the type
func R[T comparable](v T) *T {
	var zero T
	if v == zero {
		return nil
	}
	return &v
}

// NOTE: Only use this for debugging never in production code.
func DebugJSON(a any) string {
	j, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(j)
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))] //nolint:gosec
	}
	return string(b)
}
