# st

A generator for configurable types

## Background

Check [Dave Cheney on "Functional Options for Friendly APIs"](http://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis).

## How it works

This utility generates types and functions for attributes using simple rules.

Let's assume that you want a struct called `Foo`:

```go
type Foo struct {
     Addr string
     Name string
     Other int
}
```

And then, you would follow on to add Options for it:

```go
type Option func(*Foo) error

func Addr(s string) Option {
     return func(f *Foo) error {
            f.Addr = s
            return nil
     }
}

func Name(n string) Option {
     return func(f *Foo) error {
            f.Name = n
            return nil
     }
}

// etc. Boring!

func New(opts ...Option) (*Foo, error) {
     // range over opts, apply to new(Foo), etc
}
```

That initial boilerplate is what st cares about. It generates options for your struct and all the attributes, so you can focus on fine-grained Option creation using building blocks that are easy to reason about.

In order to create our `Foo` struct using st, we would simply run:

```
$ st --type Foo --fields "Name:string Addr:string Other:int" foo
```

And voilà, that initial boilerplate is created.

## Installation

I'm just getting started with this project, so be patient (also, PRs are very welcoming!).

For the time being, you have to:

- ensure you have [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports) installed.
- clone this repository
- download [gb](https://getgb.io)
  - if you have a working Go installation, run `go get github.com/constabulary/gb/...`
- cd inside the cloned repository
- build it with `gb build`
  - if you use zsh, that may conflict with the `git branch` alias, to fix it just run `unalias gb` and run `gb build` again.
- move `bin/st` somewhere on your $PATH

Subsequent functionality will improve developer experience all around, as well as code organization.