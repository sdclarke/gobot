package game

type side int

const (
	South side = iota
	North
	None
)

func (s side) String() string {
	return [...]string{"South", "North"}[s]
}

func (s side) Opposite() side {
	if s == None {
		return None
	}
	if s == South {
		return North
	}
	return South
}
