package mtk

// Three dimensional Cartesian line
type Line [2]Vector3

// Construct a Line from its endpoints
func NewLine(p, q Vector3) Line {
	return Line{p, q}
}

// Get the length
func (l Line) Length() float64 {
	return l[1].Sub(l[0]).Mag()
}

// Get the direction unit vector
func (l Line) Direction() Vector3 {
	return l[1].Sub(l[0]).Unit()
}
