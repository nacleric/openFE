package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func add(a, b int) int {
	return a + b
}

func TestAdd(t *testing.T) {
	// Given
	a := 1
	b := 2

	// When
	r := add(a, b)

	// Then
	assert.Equal(t, r, 3)
}

// [(0 0) | (0 1) | (0 2)]
// [(1 0) | (1 1) | (1 2)]
// [(2 0) | (2 1) | (2 2)]
func TestReachableCells(t *testing.T) {
	// Given
	maxMoveDistance := 1
	gridSize := 3
	unitInfo := RPG{Job: NOBLE, Movement: 2}
	u := CreateUnit(0, UnitSprite, unitInfo, PosXY{0, 1})
	i := CreateUnit(1, UnitSprite, unitInfo, PosXY{1, 0})

	units := []Unit{u, i}
	unitPointers := make([]*Unit, len(units))

	grid := CreateMGrid(unitPointers, gridSize)

	expected_legalPositions := []PosXY{{0, 1}, {2, 1}, {1, 0}, {1, 2}}

	// When
	sut := reachableCells(&grid, PosXY{1, 1}, maxMoveDistance)

	// Then
	assert.Equal(t, sut, expected_legalPositions)
}
