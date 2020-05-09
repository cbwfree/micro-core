package jwt

var (
	tokens = make(map[string]*Token)
)

func New(name string, opts ...Option) {
	tokens[name] = NewToken(opts...)
}

func Get(name string) *Token {
	return tokens[name]
}
