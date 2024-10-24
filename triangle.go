package mtk

// Three dimensional Cartesian triangle
type Triangle [3]Vector3

// Get the normal vector (not necessarily a unit vector)
func (t Triangle) Normal() Vector3 {
	pq := t[1].Sub(t[0])
	pr := t[2].Sub(t[0])
	return pq.Cross(pr)
}

// Get the unit normal vector
func (t Triangle) UnitNormal() Vector3 {
	return t.Normal().Unit()
}

// Get the area
func (t Triangle) Area() float64 {
	pq := t[1].Sub(t[0])
	pr := t[2].Sub(t[0])
	return 0.5 * pq.Cross(pr).Mag()
}
