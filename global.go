package comshim

var global = NewLoader()

func Require() error {
	return global.Load()
}

func Done() {
	global.Unload()
}

func Wait() {
	global.Wait()
}
