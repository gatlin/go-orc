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

    package main

    import (
        "fmt"
        "net/http"
        "time"
        "github.com/gatlin/go-orc"
    )

    func main() {

        // a site is a function which publishes a value asynchronously
        load_page := orc.Site{
            func(url orc.Void, publish orc.Voidchan) {
                resp, _ := http.Get(url.(string))
                defer resp.Body.Close()
                publish <- url
            },
        }

        // the built-in Cut function calls multiple sites concurrently and yields
        // the first return value
        loaded_first := orc.Cut([]orc.Voidchan{
            load_page.Call("http://archlinux.fr"),
            load_page.Call("http://google.com"),
            load_page.Call("http://amazon.com"),
        })

        // for now, it's your job to assert types.
        // Void is interface{} and Voidchan is chan interface{}
        fmt.Println(loaded_first.(string))

        // rtimer publishes `t` after `t` seconds
        rtimer := orc.Site{
            func(t orc.Void, publish orc.Voidchan) {
                <-time.After(time.Duration(t.(int)) * time.Second)
                publish <- t
            },
        }

        // this site wraps fmt.Println
        site_print := orc.Site{
            func(msg orc.Void, publish orc.Voidchan) {
                fmt.Println(msg)
                publish <- nil
            },
        }

        // sites can be recursive if pre-declared
        // metronome also shows a nifty heart-beat pattern in Orc which could be
        // useful for fault-tolerance in distributed computations
        var metronome orc.Site
        metronome = orc.Site{
            func(t orc.Void, publish orc.Voidchan) {
                orc.Merge([]orc.Voidchan{
                    site_print.Call("tick"),
                    rtimer.Call(t).ForEachDo(metronome),
                })
            },
        }

        <-metronome.Call(1)
    }

2. Acknowledgements
---

The work is inspired by the [language of the same name][1], and the work of Dr
Jayadev Misra and Dr William Cook. I am not affiliated with them.

3. Explain this Orc thing
---

Orc, is a language and concurrency calculus developed at the University of
Texas, Austin. Based on the JVM, it is designed from the ground up to
orchestrate concurrent operations and provides a nice functional syntax.

I've implemented some semantics of Orc in Go.

In Orc you have *sites* which are like functions except they "publish" values
non-deterministically and may actually be implemented across the network. On
top of this concept are four combinators: *parallel*, *sequence*, *prune*, and
*otherwise*.

*Parallel* is pretty simple: given two site invocations, it re-publishes both
of their values in tandem. `F | G` publishes the results of both `F` and `G`,
which may be sites or other expressions built on them.

*Sequence* is also simple. `F >x> G(x)` means "do F, then with the results do
G." If `F` is actually an expression (say, `H | J`), then *each* published
value of `F` is run through `G`.

*Prune* is conceptually a little harder. `F <x< G` means "Do F and G in
parallel, but hold the parts that rely on a value from G until you get one."
The first value published by the expression G is used in F.

*Otherwise* returns to simplicity: `F ; G` means "do F, and if you get a *nil*
value, publish G instead."

4. How this translates to Go
---

The actual Orc language makes certain things implicit - sites essentially are
asynchronous message queues that can return anything at any time, and the
combinators specified above really just implement schemes to manage
non-deterministic values coming out of the pipes.

Thus, sites are represented by the Site struct, which contains a function and a
`Call` method. Site functions don't return, but publish values to the channel
supplied to them. I chose to implement Sites this way because in the future I'd
like them to abstract where the Sites were defined and form the basis of
something akin to distributed objects.

*Parallel* really just merges the output channels of two site calls (which
happen concurrently thanks to `Call` using the `go` keyword), so it's
represented as `Merge`: merge a slice of void channels into a single one.

*Sequence* becomes the void channel method `ForEachDo` which accepts a Site to
be called for each value; *sequence* is essentially an implicit event loop, but
Go has explicit loops.

*Prune* takes a huge semantic hit because at the moment, without
metaprogramming, I'm not sure how I would know which parts of the
right-hand-side depend on the left - so it's assumed that all of it does. In
that regard, *prune* becomes the void channel method `WithFirstDo`. It accepts
a site just like `ForEachDo` and on the first value does something. So it's
similar in utility and semantics but not identical.

There will be practical and semantic changes; stay tuned.

5. Future
---

Right now I play fast and loose with the type system. I would eventually like
to use the reflect package to allow library users to declare normal types, and
let Orc simply ensure they're consistent. For now, though, we have type
assertions.

Another goal is already suggested by my use of Site objects rather than simple
functions: I would like for Sites to allow both local *and remote*
computations. This requires a bit more thinking but I imagine net-chans and
some kind of uniform Site registration and announcement system could make
wide-area orchestration much more feasible.

If you write a cool function building on the ones supplied here, send me a pull
request and I'll probably include it. It'd be neat to turn this into the
one-stop-shop for concurrency patterns.

Also, better examples.

[1]: http://orc.csres.utexas.edu
