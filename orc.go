package orc

/*
 * end run around the type system
 */
type Void interface{}
type Voidchan chan Void
type anyfunc func(a Void)

/*
 * Orc inspired functions
 */
func Merge(cs []Voidchan) Voidchan {
	numVals := 0
	for _, c := range cs {
		numVals = numVals + cap(c)
	}
	combined := make(Voidchan, numVals)

	for _, c := range cs {
		go func(ch Voidchan) {
			for v := range ch {
				combined <- v
			}
			return
		}(c)
	}
	return combined
}

func (self Voidchan) ForEachDo(fn anyfunc) {
	count := 0
	for v := range self {
		go fn(v)
		count = count + 1
		if count == cap(self) {
			return
		}
	}
}

func (self Voidchan) WithFirstDo(fn anyfunc) {
	v := <-self
	go fn(v)
}

func Cut(cs []Voidchan) Void {
	c := make(Voidchan, 1)
	Merge(cs).WithFirstDo(func(arg Void) {
		c <- arg
	})
	return <-c
}

/*
 * Sites, functions which remote / asynchronous services
 */
type Site struct {
	Fn func(arg Void, out Voidchan)
}

func (s Site) Call(arg Void) Voidchan {
	out := make(Voidchan, 1)
	go s.Fn(arg, out)
	return out
}
