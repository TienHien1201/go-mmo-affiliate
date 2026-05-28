package response

type responseOptions struct {
	doc  string
	meta map[string]any
	code string
}
type ResponseOption func(*responseOptions)

func defaultOptions() *responseOptions {
	return &responseOptions{
		doc:  "",
		meta: nil,
		code: "",
	}
}
func (o *responseOptions) apply(opts []ResponseOption) {
	for _, opt := range opts {
		opt(o)
	}
}

func WithDoc(doc string) ResponseOption {
    return func(o *responseOptions) {
        o.doc = doc
    }
}

func WithMeta(meta map[string]any) ResponseOption {
    return func(o *responseOptions) {
        o.meta = meta
    }
}

func WithCode(code string) ResponseOption {
    return func(o *responseOptions) {
        o.code = code
    }
}