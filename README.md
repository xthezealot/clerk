# clerk [![GoDoc](https://godoc.org/github.com/arthurwhite/clerk?status.svg)](https://godoc.org/github.com/arthurwhite/clerk) [![Build](https://travis-ci.org/arthurwhite/clerk.svg?branch=master)](https://travis-ci.org/arthurwhite/clerk) [![Coverage](https://coveralls.io/repos/github/arthurwhite/clerk/badge.svg?branch=master)](https://coveralls.io/github/arthurwhite/clerk?branch=master) [![Go Report](https://goreportcard.com/badge/github.com/arthurwhite/clerk)](https://goreportcard.com/report/github.com/arthurwhite/clerk) ![Status Stable](https://img.shields.io/badge/status-stable-brightgreen.svg)

Package [clerk](https://godoc.org/github.com/arthurwhite/clerk) provides file persistence for your Go structures.

This can replace external databases for small projects that can keep data in memory and do few writes.

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

### Init

Make a new database structure by embedding [DB](https://godoc.org/github.com/arthurwhite/clerk#DB) and initiate it with the destination file name used on saving.  
Only exported fields can be saved.

Use [DB.Rebase](https://godoc.org/github.com/arthurwhite/clerk#DB.Rebase) to retreive the data from an existent file at startup.

```Go
var db = new(struct {
	clerk.DB
	Data interface{}
})

func init() {
	clerk.Init("data.gob", db)
	if err := db.Rebase(); err != nil {
		panic(err)
	}
}
```

### Modification

[DB](https://godoc.org/github.com/arthurwhite/clerk#DB) embeds a [sync.RWMutex](https://golang.org/pkg/sync/#RWMutex) that should be used when accessing the data from multiple goroutines.

```Go
db.Lock()
defer db.Unlock()

db.Data = "one"
```

### Saving

After a data change and while the concurrent access is still locked, use [DB.Save](https://godoc.org/github.com/arthurwhite/clerk#DB.Save) to encode the source with [gob](https://golang.org/pkg/encoding/gob/) and save it in the destination file:

```Go
if err := db.Save(); err != nil {
	panic(err)
}
```
