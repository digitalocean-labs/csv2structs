package csv2structs

type Option func(*options)

type options struct {
	HeaderType      HeaderType
	HeaderTransform func(string) string
}

type HeaderType int32

const (
	HeaderTypeNone  HeaderType = iota // do not munge csv headers
	HeaderTypeSnake                   // snake_case csv headers
)

// WithHeaderType returns an Option to set the header type depending on your CSV data
func WithHeaderType(ht HeaderType) Option {
	return func(c *options) {
		c.HeaderType = ht
	}
}

// WithHeaderTransform returns an Option to set a custom header transformation function
func WithHeaderTransform(fn func(string) string) Option {
	return func(c *options) {
		c.HeaderTransform = fn
	}
}

func getOptions(opts []Option) *options {
	o := &options{
		HeaderType: HeaderTypeSnake,
	}

	for _, opt := range opts {
		opt(o)
	}

	return o
}
