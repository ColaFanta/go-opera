package opera

type Unit interface {
	IsZero() bool
	IsUnit() bool
	uniq() *struct{}
}

type unit struct {
	u *struct{}
}

func (unit) IsZero() bool {
	return true
}

func (u unit) IsUnit() bool {
	return u.u != nil
}

func (u unit) uniq() *struct{} {
	return u.u
}

var U Unit = unit{u: &struct{}{}}
