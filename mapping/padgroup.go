package mapping

import "fmt"

type PadGroup struct {
	fecID          int
	padGroupTypeID int
	padSizeID      int
	x              float64
	y              float64
}

func (pg PadGroup) String() string {
	return fmt.Sprintf("fecID %d padGroupTypeID %d padSizeID %d x %7.2f y %7.2f",
		pg.fecID, pg.padGroupTypeID, pg.padSizeID, pg.x, pg.y)
}
