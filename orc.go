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
	combined := make(Voidchan, 100)
	for _, c := range cs {
		go func(ch Voidchan) {
			for v := range ch {
				combined <- v
			}
		}(c)
	}
	return combined
}

func (self Voidchan) ForEachDo(fn anyfunc) {
	for v := range self {
		go fn(v)
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
