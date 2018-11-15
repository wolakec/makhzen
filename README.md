# Makhzen
Makhzen is a distributed in-memory data store allowing you to save string values against a key. The purpose of this project is to learn more about golang and distributed systems.

Values sent to one instance are asynchronously sent to other Makhzen instances, allowing you to retrieve the saved values from other instances. Makhzen provides AP from the CAP theorem, this means:

* Availability - the cluster is able to tolerate node failures
* Parition tolerance - the cluster can continue to function during a network partion that renders nodes unreachable

The word "makhzen" means "warehouse" in arabic (مخزن‎).

[![Build Status](https://travis-ci.com/wolakec/makhzen.svg?branch=master)](https://travis-ci.com/wolakec/makhzen)

### Limitations
- Very rudimentary and not for use in production
- Makhzen is not currently thread safe
- Makhzen can currently only store strings

### Prerequisites
This project depends upon the following:
  - Golang (https://wwww.golang.org)

### Installation
To install the package run the following command in the root directory:
```
go build main.go
```

### Running a Makhzen cluster
To run Makhzen cluster, start up each instance with the following command, each on it's own port, supplying the hosts of the other instances in the cluster argument. 
```
go run main.go -port=3000 -cluster=http://127.0.0.1:3001,http://127.0.0.1:3002
```

For example, to start up a cluster of three instances, you can run the following commands, each in a seperate terminal tab/window.
```
go run main.go -port=3001 -cluster=http://127.0.0.1:3002,http://127.0.0.1:3003
go run main.go -port=3002 -cluster=http://127.0.0.1:3003,http://127.0.0.1:3001
go run main.go -port=3003 -cluster=http://127.0.0.1:3001,http://127.0.0.1:3002
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
