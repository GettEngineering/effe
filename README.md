# Effe

[![Build Status](https://travis-ci.com/GettEngineering/effe.svg?branch=master)](https://travis-ci.com/GettEngineering/effe)
[![codecov](https://codecov.io/gh/GettEngineering/effe/branch/master/graph/badge.svg)](https://codecov.io/gh/GettEngineering/effe)
[![codebeat badge](https://codebeat.co/badges/9c74d700-ebf8-4b76-8405-1950874576c4)](https://codebeat.co/projects/github-com-maratori-testpackage-master)
[![Maintainability](https://api.codeclimate.com/v1/badges/bf753d7560c8e4aa5cf0/maintainability)](https://codeclimate.com/github/GettEngineering/effe/maintainability)
[![Go Report Card](https://goreportcard.com/badge/github.com/GettEngineering/effe)](https://goreportcard.com/report/github.com/GettEngineering/effe)
[![GitHub](https://img.shields.io/github/license/GettEngineering/effe.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/GettEngineering/effe?status.svg)](http://pkg.go.dev/github.com/GettEngineering/effe)

**Effe** is an _orchestration_ engine for business-logic

## Installing

```go
go get github.com/GettEngineering/effe/cmd/effe
```

## Run

```bash
$ effe -h
Usage of effe:
  -d    draw diagrams for business flows
  -out string
        draw output directory (default "graphs")
  -v    show current version of effe
```

and ensuring that `$GOPATH/bin` is added to your `$PATH`.

## Documentation & Getting Started

http://gettengineering.github.io/effe

[Getting Started](https://gettengineering.github.io/effe/gettingstarted/basicconcepts/) guide.

## Issues/Problems/Ideas

https://github.com/GettEngineering/effe/issues

## Get support

Effe is maintained by Gett. Use github issue tracking for any support request.

## License

Copyright (c) 2020 Gett

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
