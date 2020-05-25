package object

type Env struct {
	store map[string]Object
	outer *Env
}

func NewEnv(outer *Env) Env {
	return Env{
		store: make(map[string]Object),
		outer: outer,
	}
}

func (e *Env) Get(id string) (Object, bool) {
	obj, ok := e.store[id]
	if !ok && e.outer != nil {
		return e.outer.Get(id)
	}
	return obj, ok
}

func (e *Env) Set(id string, o Object) {
	e.store[id] = o
}
