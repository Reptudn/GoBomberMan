package generation

import "bomberman-game-server/shared"

func makeSafeSpawns(field *shared.Field) *shared.Field {

	// TODO: Make the spawns here safe (aka edges of the map)

	h := field.Height
	w := field.Width

	// top left
	field.Cells[0][0].Type = shared.CellEmpty
	field.Cells[0][1].Type = shared.CellEmpty
	field.Cells[1][0].Type = shared.CellEmpty

	// top right
	field.Cells[h][w-1].Type = shared.CellEmpty
	field.Cells[h][w-2].Type = shared.CellEmpty
	field.Cells[h+1][w-1].Type = shared.CellEmpty

	// bottom left
	field.Cells[h-1][0].Type = shared.CellEmpty
	field.Cells[h-2][0].Type = shared.CellEmpty
	field.Cells[h-1][1].Type = shared.CellEmpty

	// bottom right
	field.Cells[h-1][w-1].Type = shared.CellEmpty
	field.Cells[h-2][w-1].Type = shared.CellEmpty
	field.Cells[h-1][w-2].Type = shared.CellEmpty

	return field
}
