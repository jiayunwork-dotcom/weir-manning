package cross

type Shape int

const (
	ShapeRect Shape = iota
	ShapeTrapezoid
)

type Section struct {
	Shape Shape
	B     float64
	M     float64
}

func NewRect(b float64) (Section, error) {
	return NewSection(ShapeRect, b, 0)
}

func NewTrapezoid(b, m float64) (Section, error) {
	return NewSection(ShapeTrapezoid, b, m)
}

func NewSection(shape Shape, b, m float64) (Section, error) {
	s := Section{Shape: shape, B: b, M: m}
	if err := ValidateSection(s); err != nil {
		return Section{}, err
	}
	return s, nil
}

func (s Section) IsRect() bool {
	return s.Shape == ShapeRect
}

func (s Section) IsTrapezoid() bool {
	return s.Shape == ShapeTrapezoid
}

func (s Section) Label() string {
	if s.IsRect() {
		return "rect"
	}
	return "trapz"
}
