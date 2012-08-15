Orc
===

Concurrency orchestration patterns for a concurrent language

(c) 2012 Gatlin Johnson <rokenrol@gmail.com>. Licensed under the WTFPL.

0. Introduction
---

Orc is a simple library to aid with concurrency orchestration. The goal is to
write useful patterns once and to build more complex patterns out of
intelligent, simple primitives.

1. Synopsis
---

    func main() {
        // first example: print which site loads first
        // this site loads a URL and then publishes the URL when it finishes
        s1 := Site{
            func(url Void, out Voidchan) {
                resp, _ := http.Get(url.(string))
                defer resp.Body.Close()
                out <- url
            },
        }

        res1 := Cut([]Voidchan{
            s1.Call("http://archlinux.fr"),
            s1.Call("http://easynews.com"),
            s1.Call("http://google.com"),
        })
        fmt.Println(res1.(string))

        // second example: interleave the ongoing results of two concurrent
        // operations
        // this site publishes "m" every "m" seconds

        s2 := Site{
            func(m Void, out Voidchan) {
                for {
                    <-time.After(time.Duration(m.(int)) * time.Second)
                    out <- strconv.Itoa(m.(int))
                }
            },
        }
        res2 := Merge([]Voidchan{s2.Call(2), s2.Call(3)})
        res2.ForEachDo(func(s Void) {
            fmt.Println(s)
        })
    }

2. Acknowledgements
---

The work is inspired by the [language of the same name][1], and the work of Dr
Jayadev Misra and Dr William Cook. I am not affiliated with them.

3. Explain this Orc thing
---

The Orc language, as a domain specific language built from scratch, makes
certain language constructs implicit: asynchronous queues, iteration, and even
concurrency itself. The corner stone of their work is the *site*, which is like
a function except it executes concurrently and *publishes* an arbitrary number
of values at non-deterministic times.

In Go, we have functions, the `go` keyword, and asynchronous queues as
first-class values. Thus, a *site* is any function which returns a chan
(specifically, type `chan interface{}` or `Voidchan` for convenience).

On top of this notion of sites there are three combinators: *parallel*,
*sequence*, and *prune*. These ideas have been modified to reflect the
mechanics of the host language while retaining (hopefully) the same semantics.

*Parallel* takes two site calls and re-publishes their independent results 
together. Since really all this is doing is merging the output streams of two
chans (Go handles the concurrency part), this library provides `Merge`.

*Sequence* takes a source of published values on the left (such as a site call,
or an application of *Parallel*) and an expression on the right using the
publish values. Since this is basically an iteration over a channel, this
library supplies the method `ForEachDo`, which accepts a lambda to process the
values.

Finally, *Prune* takes a source of published values and returns the first one.
In the spirit of `ForEachDo`, this library supplies `WithFirstDo`. It has the
same signature.

Together, `Merge`, `ForEachDo`, and `WithFirstDo` allow us to build highly
concurrent applications and express common patterns elegantly, all the while
retaining guarantees about the effects of our programs.

Orc additionally intends on supplying common combinations of these functions;
currently, `Cut` is also provided (a combination of `Merge` and `WithFirstDo`).

4. Future
---

Right now I play fast and loose with the type system. I would eventually like
to use the reflect package to allow library users to declare normal types, and
let Orc simply ensure they're consistent. For now, though, we have type
assertions.

If you write a cool function building on the ones supplied here, send me a pull
request and I'll probably include it. It'd be neat to turn this into the
one-stop-shop for concurrency patterns.

Also, better examples.

[1]: http://orc.csres.utexas.edu
