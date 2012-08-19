package orc

/*
 * end run around the type system
 */

type Void interface{}
type Voidchan chan Void

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

/*
 * Orc inspired functions
 */

func Merge(cs []Voidchan) Voidchan {
	combined := make(Voidchan, cap(cs))
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

func (self Voidchan) ForEachDo(site Site) Voidchan {
	out := make(Voidchan, cap(self))
	for i := 0; i < cap(self); i++ {
		s := <-self
		out <- site.Call(s)
	}
	return out
}

func (self Voidchan) WithFirstDo(fn Site) {
	v := <-self
	<-fn.Call(v)
}

func Cut(cs []Voidchan) Void {
	c := make(Voidchan, 1)
	Merge(cs).WithFirstDo(Site{
		func(arg Void, out Voidchan) {
			c <- arg
			out <- nil
		},
	})
	return <-c
}
