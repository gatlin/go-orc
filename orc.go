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
	numSites := len(cs)
	combined := make(Voidchan, numSites+1)
	halter := make(chan int, numSites)
	go func(h chan int) {
		i := 0
		for n := range h {
			i = i + n
			if i == numSites {
				combined <- nil
				return
			}
		}
	}(halter)

	for _, c := range cs {
		go func(ch Voidchan, h chan int) {
			for v := range ch {
				combined <- v
				h <- 1
			}
			return
		}(c, halter)
	}
	return combined
}

func (self Voidchan) ForEachDo(fn anyfunc) {
	for v := range self {
		if v == nil {
			return
		} else {
			go fn(v)
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
	out := make(Voidchan, 100)
	go s.Fn(arg, out)
	return out
}
