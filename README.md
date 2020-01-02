# elasthink
An alternative to elasticsearch engine for small set of documents that makes you think a little bit when implementing it.
We use inverted index to build the index. We also utilize redis to store the indexes.

## Table of Contents

* [Installation](#installation)
* [Documentation](#documentation)
* [Dependencies](#dependencies)

## Installation
1. To install elasthink, you need to run `$ go get github.com/SurgicalSteel/elasthink`
2. Then you need to specify your redis addresses for each environment in `files/config/redis` folder
3. To start with your own document, you need to modify the document type const in `entity/document.go` and its validation function in `module/document.go`
4. Everything is set, you can now build and run elasthink.

## Documentation
1. API documentation (insomnia format) is available in the `elasthink_insomnia_api_documentation.json`
2. To view code documentation, you can visit [Elasthink's GoDoc](https://godoc.org/github.com/SurgicalSteel/elasthink)

## dependencies
1. [gorilla/mux]("https://github.com/gorilla/mux")
2. [gomodule/redigo](https://github.com/gomodule/redigo)
