package core

import (
	"fmt"
	"testing"
)

// [1,2,3]
// [4,5,6]
// [7,8,9]
func reachableCells_test(t *testing.T) {
	// Given
	maxMoveDistance := 1
	gridSize := 3
	unitInfo := RPG{Job: NOBLE, Movement: 2}
	u := CreateUnit(0, UnitSprite, unitInfo, PosXY{0, 1})
	i := CreateUnit(1, UnitSprite, unitInfo, PosXY{1, 0})

	units := []Unit{u, i}
	unitPointers := make([]*Unit, len(units))

	grid := CreateMGrid(unitPointers, gridSize)
	fmt.Println(grid)

	// When
	sut := reachableCells(&grid, PosXY{1, 1}, maxMoveDistance)

	// Then
}
