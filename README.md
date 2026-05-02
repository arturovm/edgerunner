# edgerunner

A tool for running concurrent workloads.

## Why

I got tired of writing very similar tools for specific use cases while solving
challenges on _picoCTF_, so I decided to abstract the pattern, and I wrote this.

## How

### Installing it

If your Go environment is setup properly, you can run:

```sh
$ go install git.sr.ht/~arturovm/edgerunner@latest
```

No pre-compiled binaries are available as of now.

### Using it

The concepts to understand when using `edgerunner` are the _generator_ and
_task_:

1. The _generator_ is a function that returns a list of inputs that will later
   be distributed to the concurrent workers executing the _task_.
2. The _task_ is a function that takes as its single input each value produced
   by the _generator_ and then does something useful with it (e.g. using the
   values as URIs to enumerate paths on a server, etc).

#### Command usage

```
Usage of edgerunner:
  -g, --generator string   path of generator spec file (default "generator.fnl")
  -t, --task string        path of task spec file (default "task.fnl")
  -w, --workers int        number of task runners
```

`edgerunner` defaults to spawning a number of workers equal to the number of
CPUs on your system, but you can of course modify this behavior with the
flag provided for this purpose.

#### Writing a generator and task spec

`edgerunner` uses [Fennel](https://fennel-lang.org) as its scripting language
for both the _generator_ and the _task_ specs. Why? Because I love Functional
Programming, I like LISPs, and because I can :)

In this repository, you'll find some examples of a minimal _generator_ and _task_
under `examples`, but they are reproduced here for your benefit:

```fennel
;; generator.fnl
(fn generator []
  [1 2 3])
```

```fennel
;; task.fnl
(local crypto (require :crypto))

(fn task [val]
  (print (crypto.md5 val)))
```

The most important thing to remember here is that the `generator` and `task`
functions _must_ be the last expressions in the spec files. This is due to the
underlying mechanics of Fennel evaluation.

`edgerunner` exposes a standard library of Go functions, thanks to the power of
[`gopher-lua-libs`](https://github.com/vadv/gopher-lua-libs). Please review the
documentation of that package to understand what's available.

### Building it

Clone the repository and run:

```sh
$ make build
```

## License

`edgerunner` is distributed under the MIT License. You can find more details in
`LICENSE`.

## Important

You won't find any AI slop here.
