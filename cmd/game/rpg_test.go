package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	TEST_GRIDSIZE = 3
)

func TestAdd(t *testing.T) {
	// Given
	x := 1
	y := 2
	expected_res := 3
	// When
	res := Add(x, y)

	// Then
	assert.Equal(t, expected_res, res)
}
