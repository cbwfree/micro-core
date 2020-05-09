package event

var (
	DefaultEvents = NewEvents("default")
)

func Exist(name string) bool {
	return DefaultEvents.Exist(name)
}

func Bind(name string, handler Handler, maxCount ...int64) error {
	return DefaultEvents.Bind(name, handler, maxCount...)
}

func Once(name string, handler Handler) error {
	return DefaultEvents.Once(name, handler)
}

func Delete(name string) {
	DefaultEvents.Delete(name)
}

func Call(name string, obj interface{}) error {
	return DefaultEvents.Call(name, obj)
}
