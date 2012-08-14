Orc
===

Concurrency orchestration patterns for a concurrent language

(c) 2012 Gatlin Johnson <rokenrol@gmail.com>

0. Introduction
---

Orc is a simple library to aid with concurrency orchestration. The goal is to
write useful patterns once and to build more complex patterns out of
intelligent, simple primitives.

1. Synopsis
---

    // which site loads first?

    func fetch_site(url string) Voidchan {
        out := make(Voidchan, 1)
        go func() {
            resp, _ := http.Get(url)
            defer resp.Body.Close()
            out <- url
        }()
        return out
    }

    func main() {
        res := orc.Cut([]Voidchan{
            fetch_site("http://archlinux.fr"),
            fetch_site("http://www.easynews.com"),
            fetch_site("http://golang.org"),
        })
        fmt.Printf("%s loaded first\n", res1.(string))
    }

2. Acknowledgements
---

The work is inspired by the [language of the same name][1], and the work of Dr
Jayadev Misra and Dr William Cook.

3. This doesn't look much like Orc ... / can you explain this?
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

[1]: http://orc.csres.utexas.edu
