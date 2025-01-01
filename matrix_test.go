package mtk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatrixSetValue(t *testing.T) {
	m := NewMatrix(3, 4)

	shape := m.Shape()

	assert.Equal(t, 3, shape[0])
	assert.Equal(t, 4, shape[1])
	assert.Equal(t, 0., m.At(2, 3))
}

func TestMatrixQR(t *testing.T) {
	m := NewMatrix(3, 2)
	m.SetValue(0, 0, 1)
	m.SetValue(0, 1, 2)
	m.SetValue(1, 0, 3)
	m.SetValue(1, 1, 4)
	m.SetValue(2, 0, 5)
	m.SetValue(2, 1, 6)

	q, r := m.QR()

	assert.Equal(t, 3, q.Shape()[0])
	assert.Equal(t, 2, q.Shape()[1])
	assert.Equal(t, "0.169", fmt.Sprintf("%.3f", q.At(0, 0)))
	assert.Equal(t, "0.897", fmt.Sprintf("%.3f", q.At(0, 1)))
	assert.Equal(t, "0.507", fmt.Sprintf("%.3f", q.At(1, 0)))
	assert.Equal(t, "0.276", fmt.Sprintf("%.3f", q.At(1, 1)))
	assert.Equal(t, "0.845", fmt.Sprintf("%.3f", q.At(2, 0)))
	assert.Equal(t, "-0.345", fmt.Sprintf("%.3f", q.At(2, 1)))

	assert.Equal(t, 2, r.Shape()[0])
	assert.Equal(t, 2, r.Shape()[1])
	assert.Equal(t, "5.916", fmt.Sprintf("%.3f", r.At(0, 0)))
	assert.Equal(t, "7.437", fmt.Sprintf("%.3f", r.At(0, 1)))
	assert.Equal(t, "0.000", fmt.Sprintf("%.3f", r.At(1, 0)))
	assert.Equal(t, "0.828", fmt.Sprintf("%.3f", r.At(1, 1)))
}

func TestMatrixSymmetricEigen(t *testing.T) {
	m := NewMatrix(3, 3)
	m.data = []float64{4, 5, -1, 5, 0, -7, -1, -7, 9}

	e, v := m.SymmetricEigen()

	assert.Equal(t, "13.995", fmt.Sprintf("%.3f", e[0]))
	assert.Equal(t, "-5.530", fmt.Sprintf("%.3f", e[1]))
	assert.Equal(t, "4.535", fmt.Sprintf("%.3f", e[2]))

	assert.Equal(t, "0.336", fmt.Sprintf("%.3f", v.At(0, 0)))
	assert.Equal(t, "0.515", fmt.Sprintf("%.3f", v.At(1, 0)))
	assert.Equal(t, "-0.789", fmt.Sprintf("%.3f", v.At(2, 0)))
	assert.Equal(t, "-0.399", fmt.Sprintf("%.3f", v.At(0, 1)))
	assert.Equal(t, "0.836", fmt.Sprintf("%.3f", v.At(1, 1)))
	assert.Equal(t, "0.375", fmt.Sprintf("%.3f", v.At(2, 1)))
	assert.Equal(t, "0.853", fmt.Sprintf("%.3f", v.At(0, 2)))
	assert.Equal(t, "0.189", fmt.Sprintf("%.3f", v.At(1, 2)))
	assert.Equal(t, "0.487", fmt.Sprintf("%.3f", v.At(2, 2)))
}
