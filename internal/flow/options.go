package flow

type SolveOptions struct {
	Tolerance       float64
	MaxIterations   int
	InitialHigh     float64
	BracketExpand   float64
	LowDepth        float64
	MaxBracketSteps int
	Ceiling         float64
}

func DefaultOptions() SolveOptions {
	return SolveOptions{
		Tolerance:       1e-8,
		MaxIterations:   200,
		InitialHigh:     1.0,
		BracketExpand:   2.0,
		LowDepth:        1e-6,
		MaxBracketSteps: 40,
		Ceiling:         1e9,
	}
}

func (o SolveOptions) WithTolerance(t float64) SolveOptions {
	o.Tolerance = t
	return o
}

func (o SolveOptions) WithMaxIterations(n int) SolveOptions {
	o.MaxIterations = n
	return o
}

func (o SolveOptions) Validate() error {
	if o.Tolerance <= 0 {
		return ErrTolerance
	}
	if o.MaxIterations <= 0 {
		return ErrIterations
	}
	if o.InitialHigh <= 0 {
		return ErrBracket
	}
	if o.BracketExpand <= 1 {
		return ErrBracket
	}
	if o.LowDepth <= 0 {
		return ErrBracket
	}
	if o.Ceiling <= o.InitialHigh {
		return ErrCeiling
	}
	return nil
}

var (
	ErrTolerance   = errorOf("tolerance must be positive")
	ErrIterations  = errorOf("max iterations must be positive")
	ErrBracket     = errorOf("bracket parameters are invalid")
	ErrCeiling     = errorOf("depth ceiling must exceed the initial high bound")
)

func errorOf(msg string) error {
	return &optionsError{msg: msg}
}

type optionsError struct {
	msg string
}

func (e *optionsError) Error() string {
	return e.msg
}
