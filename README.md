# clerk [![GoDoc](https://godoc.org/github.com/arthurwhite/clerk?status.svg)](https://godoc.org/github.com/arthurwhite/clerk) [![Build](https://travis-ci.org/arthurwhite/clerk.svg?branch=master)](https://travis-ci.org/arthurwhite/clerk) [![Coverage](https://coveralls.io/repos/github/arthurwhite/clerk/badge.svg?branch=master)](https://coveralls.io/github/arthurwhite/clerk?branch=master) [![Go Report](https://goreportcard.com/badge/github.com/arthurwhite/clerk)](https://goreportcard.com/report/github.com/arthurwhite/clerk) ![Status Testing](https://img.shields.io/badge/status-testing-orange.svg)

Package [clerk](https://godoc.org/github.com/arthurwhite/clerk) provides file persistence for your Go variables.

It can replace external databases for small projects that can keep data in memory and don't make a lot of writes.

## Installing

1. Get package:

	```Shell
	go get -u github.com/arthurwhite/clerk
	```

2. Import it in your code:

	```Go
	import "github.com/arthurwhite/clerk"
	```

## Usage

Make a new database with a source and a destination file name:

```Go
var data interface{}

db, err := clerk.New("data.db", &data)
if err != nil {
	panic(err)
}
```

After a data change, use [DB.Save](https://godoc.org/github.com/arthurwhite/clerk#DB.Save) to encode the source with [gob](https://golang.org/pkg/encoding/gob/) and save it in the destination file:

```Go
data = "one"
if err = data.Save(); err != nil {
	panic(err)
}
```

On restart, the data file will be decoded back in the source variable.
