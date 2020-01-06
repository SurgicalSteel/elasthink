# elasthink
An alternative to elasticsearch engine written in Go for small set of documents that makes you think a little bit when implementing it.
We use inverted index to build the index. We also utilize redis to store the indexes.

## Table of Contents

* [Installation](#installation)
* [Documentation](#documentation)
* [Dependencies](#dependencies)
* [Reference](#reference)

## Installation
1. To install elasthink, you need to run `$ go get github.com/SurgicalSteel/elasthink`
2. Then you need to specify your redis addresses for each environment in `files/config/redis` folder
3. To start with your own document, you need to modify the document type const in `entity/document.go` and its validation function in `module/document.go`
4. To build elasthink, run `$ go build`
5. To run elasthink, run `$ ./elasthink -env={your-environment}` and your elasthink web service should run on `localhost:9000` 

## Documentation
API documentation (insomnia format) is available in the `elasthink_insomnia_api_documentation.json`  
For code documentation, we use standard godoc as our code documentation tool. To view the code documentation, follow these steps :
1. Open your terminal, head to this cloned repo (SurgicalSteel/elasthink)
2. run `godoc -http=:6060` (this will trigger godoc at port 6060)
3. Open your browser, and hit `http://127.0.0.1:6060/pkg/github.com/SurgicalSteel/elasthink/`

## Dependencies
1. [gorilla/mux](https://github.com/gorilla/mux)
2. [gomodule/redigo](https://github.com/gomodule/redigo)
3. [gcfg.v1](https://gopkg.in/gcfg.v1)

## Reference
[E-Book Redis in Action Part 2 Chapter 7](https://redislabs.com/ebook/part-2-core-concepts/chapter-7-search-based-applications/7-1-searching-in-redis/7-1-1-basic-search-theory/)
