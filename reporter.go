package monitor

type Reporter interface {
	Update(float64)
}
