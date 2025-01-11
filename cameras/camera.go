package cameras

type Camera struct {
	X float64
	Y float64
}

func NewCamera(x, y float64) *Camera {
	return &Camera{
		X: x,
		Y: y,
	}
}
