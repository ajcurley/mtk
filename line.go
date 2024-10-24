package mtk

// Three dimensional Cartesian line
type Line [2]Vector3

// Get the length
func (l Line) Length() float64 {
	return l[1].Sub(l[0]).Mag()
}

// Get the direction unit vector
func (l Line) Direction() Vector3 {
	return l[1].Sub(l[0]).Unit()
}
