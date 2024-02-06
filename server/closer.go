package server

type closer []func()

// add adds a closer function.
func (c *closer) add(f func()) {
	*c = append(*c, f)
}

// closeAll closes all the closer functions.
func (c *closer) closeAll() {
	for _, f := range *c {
		f()
	}
}
