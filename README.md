# static [![Build Status](https://travis-ci.org/leighmcculloch/static.svg?branch=master)](https://travis-ci.org/leighmcculloch/static) [![Go docs](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/leighmcculloch/static)

A package for generating static websites from any Go web app that uses `net/http`.

## Why

Building static websites with existing frameworks like [middleman](https://github.com/middleman/middleman) is easy, but large websites can be slow. [hugo](https://github.com/spf13/hugo) is a popular option if you need to convert content using templates, but not if you have a go web app that you want to make static.

[static](https://github.com/leighmcculloch/static) helps you use build static websites that are dynamically generated from sources like RSS feeds, databases, APIs, etc by calling each handler registered and saving the output as files.

## Go docs

Get the go docs at: [godoc.org/github.com/leighmcculloch/static](https://godoc.org/github.com/leighmcculloch/static)

## Install

```bash
go get github.com/leighmcculloch/static
```

```go
import "github.com/leighmcculloch/static"
```

## Usage

Call `Build` with a `http.Handler`, a `[]string` of paths to build to static files, and a callback for printing progress and errors which are communicated via events. The event handler can be `nil` but it's the only way you'll find out if there's an error building a path.

```go
options := static.DefaultOptions
static.Build(options, handler, paths, func (e static.Event) {
  log.Println(e)
})
```

## Options

Instead of using the default `Options` you can define your own.

```go
options := static.Options{
  OutputDir:   "build",
  Concurrency: 50,
  DirFilename: "index.html",
}
static.Build(options, handler, paths, func (e static.Event) {
  log.Println(e)
})
```

## Simple Example

Fire up the sample below. Running the Hello World web server is as you'd expect `go run *.go`, and then building the static version is as simple as `go run *.go -build`.

```go
package main

import (
  "net/http"
  "github.com/leighmcculloch/static"
)

var build bool

func init() {
  flag.BoolVar(&build, "build", false, "Build the website to static files rather than run the web server.")
  flag.Parse()
}

func main() {
  handler := http.NewServeMux()
  paths := []string{}

  paths = append(paths, "/")
  handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Requests) {
    fmt.Fprintf(w, "Hello %s!", r.URL.Path)
  })

  if build {
    options := static.DefaultOptions
    static.Build(options, handler, paths, func (e static.Event) {
      log.Println(e)
    })
  } else {
    s := &http.Server{Addr: ":8080", Handler: handler}
    log.Fatal(s.ListenAndServe())
  }
}
```

## Typical Example

See [github.com/leighmcculloch/readprayrepeat.com](https://github.com/leighmcculloch/readprayrepeat.com).
