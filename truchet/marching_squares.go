package truchet

// TODO: get back to this one

// type MarchingSquares struct {
// 	nX    int
// 	nY    int
// 	cells map[cellCoord]*Cell
// 	// edge containers, specifying the position of cell border points
// 	// columnEdges map[cellCoord]Edge
// 	// rowEdges    map[cellCoord]Edge
// 	source samplers.DataSource
// }

// func NewMarchingGrid(b box.Box, nx int, source samplers.DataSource, threshold float64) MarchingSquares {
// 	boxes := b.PartitionIntoSquares(nx)
// 	cells := make(map[cellCoord]*Cell, len(boxes))
// 	if len(boxes) != nx*nx {
// 		panic(fmt.Errorf("not right, want %d, got %d", nx*nx, len(boxes)))
// 	}
// 	// horPoints := tileSet.EdgePointMapping.getHorizontal()
// 	// vertPoints := tileSet.EdgePointMapping.getVertical()

// 	// getIntersects := getSourcedIntersects
// 	// grid.rowEdges = make(map[cellCoord]Edge, nx+1)
// 	// for i := range nx + 1 { // for each of horizontal edges
// 	// 	for j := range nx { // for each cell
// 	// 		intersects := getIntersects(horPoints, edgeSource, float64(j)/float64(nx+1), float64(i)/float64(nx+1))
// 	// 		grid.rowEdges[cellCoord{j, i}] = Edge{intersects} // flipped order is intentional
// 	// 	}
// 	// }
// 	// grid.columnEdges = make(map[cellCoord]Edge, nx+1)
// 	// for i := range nx + 1 { // for each of vertical edges
// 	// 	for j := range nx { // for each cell
// 	// 		intersects := getIntersects(vertPoints, edgeSource, float64(i)/float64(nx+1), float64(j)/float64(nx+1))
// 	// 		grid.columnEdges[cellCoord{i, j}] = Edge{intersects}
// 	// 	}
// 	// }
// 	for _, childBox := range boxes {
// 		cell := &Cell{
// 			Grid: grid,
// 			Box:  childBox.Box,
// 			x:    childBox.I,
// 			y:    childBox.J,
// 		}
// 		cell.PopulateCurves(tilePicker)
// 		cells[cellCoord{childBox.I, childBox.J}] = cell
// 	}
// 	grid.nX = nx
// 	grid.nY = nx
// 	grid.cells = cells
// 	return grid
// }
