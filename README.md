# Makhzen
Makhzen is an in-memory data store allowing you to save string values against a key. This is a project to help me to learn Golang and hopefully play around with some interesting things (like distributed computing)

The word "makhzen" means "warehouse" in arabic (مخزن‎).

[![Build Status](https://travis-ci.com/wolakec/makhzen.svg?branch=master)](https://travis-ci.com/wolakec/makhzen)

### Prerequisites
This project depends upon the following:
  - Golang (https://wwww.golang.org)

### How to install
To install the package run the following command:
```
go get github.com/wolakec/makhzen/store
```

### Usage
```go
package main

import (
	"fmt"
	"log"

	"github.com/wolakec/makhzen/store"
)

func main() {
	s := store.New()
	s.Set("testKey", "1234")

	if item, exists := s.Get("testKey"); exists {
		fmt.Println(item)
	}

	log.Fatalf("Item does not exist")
}
```

### Limitations
- Makhzen is not currently thread safe
- Makhzen can currently only store strings
