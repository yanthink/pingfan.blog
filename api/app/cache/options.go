package cache

type Option func(o *Options)

type Options struct {
	Prefix string
}

func ApplyOptions(opts ...Option) *Options {
	o := &Options{
		Prefix: "",
	}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

func WithPrefix(prefix string) Option {
	return func(o *Options) {
		o.Prefix = prefix
	}
}
