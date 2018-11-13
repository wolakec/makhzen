# Makhzen
Makhzen is an in-memory data store allowing you to save string values against a key. This is a project to help me to learn Golang and hopefully play around with some interesting things (like distributed computing)

The word "makhzen" means "warehouse" in arabic (مخزن‎).

[![Build Status](https://travis-ci.com/wolakec/makhzen.svg?branch=master)](https://travis-ci.com/wolakec/makhzen)

### Prerequisites
This project depends upon the following:
  - Golang (https://wwww.golang.org)

### Installation
To install the package run the following command in the root directory:
```
go build main.go
```

### Running Makzen
To run Makzen run the following command, optionally supplying a port argument. The default port is 5000.
```
go run main.go -port=3000
```

### Setting a value
To set a value make a PUT request to /items, supplying the desired key in the URL.

```
curl -d '{"value":"eu-west-1"}' -H "Content-Type: application/json" -X PUT http://localhost:3000/items/region
```

### Fetching a value
To fetch a saved value make a GET request to /items, supplying the key in the URL

```
curl http://localhost:3000/items/region
```

### Response
The response is currently being sent as simple text (Will change this to JSON at some point)

### Limitations
- Makhzen is not currently thread safe
- Makhzen can currently only store strings
