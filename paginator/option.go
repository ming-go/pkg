package paginator

type Option interface {
	Apply(*paginator)
}

func WithCopySlice(b bool) Option {
	return withCopySlice(b)
}

type withCopySlice bool

func (w withCopySlice) Apply(p *paginator) {
	p.copySlice = bool(w)
}
