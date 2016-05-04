# go-raml
[![Build Status](https://travis-ci.org/Jumpscale/go-raml.svg?branch=master)](https://travis-ci.org/Jumpscale/go-raml)

## What is go-raml

When creating and maintaining api's, there are two approaches:
* Design first

    You design the api, all methods, descriptions and types from the api point of view. Afterwards, you generate all the boilerplate code and documentation to bootstrap the code.

* Code first

    When modifying the implementation, the interface definitions need to be kept in sync. Besides the good practice of keeping your specification up to date, it is used by other tools to generate clients, documentation, ...

This tool supports both (or at least, this is on the roadmap).
As a specification format, it uses [RAML 1.0](http://raml.org) .

## RAML versions
Only RAML version 1.0 RC is supported.

Currently there are still some [limitations](docs/limitations.md) on the RAML 1.0 features that are supported.

## Install

make sure you have at least go 1.6 installed !

Install `godep` as package manager

`$go get -u github.com/tools/godep`

Install go-bindata, we need it to compile the template files to .go file

`go get -u github.com/jteeuwen/go-bindata/...`

Install go-raml

`go get -u github.com/Jumpscale/go-raml`

### Build

```
cd $GOPATH/src/github.com/Jumpscale/go-raml
sh build.sh
```

## Usage

To use it on the commandline, just execute `go-raml` without any arguments, it will output the help on the stdout.

To invoke the codegeneration using `go generate`, specify the generation in 1 of your go source files:
`//go:generate go-raml ...`

go-raml needs to be on the path for this to work off course.

## Code generation

Internally, go templates are used to generate the code, this provides a flexible way to alter the generated code and to add different languages for the client.

## Server

`go-raml` able to generate Go & Python server. API Docs also generated in `/apidocs/` endpoint. API Docs generation can be disabled by specifying `no-apidocs` option.

### Go Server

To generate the Go code for implementing the server in first design approach, execute

`go-raml server -l go --dir ./result_directory --ramlfile api.raml`

The generated server uses [Gorilla Mux](http://www.gorillatoolkit.org/pkg/mux) as HTTP request multiplexer.

-Generated codestructure:
-* Interfaces types, always regenerated
-* Implementing types, only generated when the file is not present


### Flask/Python Server

To generate the Flask/Python code for implementing the server in first design approach, execute

`go-raml server -l python --dir ./result_directory --ramlfile api.raml`

### Code Generator Options

```
   --language, -l "go"  Language to construct a server for
   --dir "."        target directory
   --ramlfile "."   Source raml file
   --package "main" package name
   --no-main        Do not generate a main.go file
   --no-apidocs     Do not generate API Docs in /apidocs/ endpoint
```

## Client

`go-raml client --language go  --dir ./result_directory --ramlfile api.raml`

A go 1.5.x compatible client is generated in result_directory directory.

`go-raml client --language python --dir ./result_directory --ramlfile api.raml`

A python 3.5 compatible client is generated in result_directory directory.

## Type

[types](http://docs.raml.org/specs/1.0/#raml-10-spec-types) is mapped to struct.

Some rules about properties naming:

- capitalizing first character of the properties.
- json tag is the same as property name

File name is the same as types name with lowercase letter.
struct name = types name.

#### Type Mapping
    Raml        |  Go   
    ----------- | -----------
    string      | string
    number      | float
    integer     | int
    boolean     | bool
    date        | Date
    enum        | []string
    file        | string
    sometype[]  | []sometype
    sometype[][]| [][]sometype
    Union       | interface{}

## Bodies
[Request Body](http://docs.raml.org/specs/1.0/#raml-10-spec-bodies) and response body are mapped into structs
and following the same rules as types.

struct name = [Resource name][Method name][ReqBody|RespBody].

RequestBody generated from body node below method.

ResponseBody generated from body node below responses.

## Resource
[Resource](http://docs.raml.org/specs/1.0/#raml-10-spec-resources-and-nested-resources) in the server is mapped to:
- interface:

    - file name = [resource]_if.go
    - always be regenerated
    - interface name = [resource]Interface

- API implementation that implements the interface
    - file name = [resource]_api.go
    - only generated when the file is not present
    - struct name = [resource]API

- routes for all necessary routes:
    - func name = [Resource]InterfaceRoutes


## Header

Code related to requests [headers](http://docs.raml.org/specs/1.0/#raml-10-spec-headers) are only generated in the Client lib. All functions have arguments to send any request header, the current client lib will not check the headers against the RAML specifications.


Response headers related code is only generated in the server in the form of commented code, example:
```
// uncomment below line to add header
// w.Header.Set("key","value")
```

## Query Strings and Query Parameters

All client library functions have arguments to send [Query Strings and Query Parameters](http://docs.raml.org/specs/1.0/#raml-10-spec-query-strings-and-query-parameters), the current client lib will not check it against the RAML specifications.

The generated code in the server is in the form of commented code:

```
// name := req.FormValue("name")
```

## Input Validation


    Validation      |    Go     | Python
------------------------------- | ------------- | -----------
 minLength                      |   v   |   v
 maxLength          |   v   |   v
 pattern            |   v   |   v
 minimum            |   v   |   v
 maximum            |   v   |   v
 format             |       x   |   x
 multipleOf         |   v   |   v
 array field minItems       |   v   |   v
 array field maxItems       |   v   |   v
 array field uniqueItems    |   v   |   x
 array Type minItems        |   v   |   x
 array Type maxItems        |   v   |   x
 array Type uniqueItems     |   v   |   x

## Specification file

Besides generation of a new RAML specification file, updating an existing raml file is also supported. This way the raml filestructure that can be included in the main raml file is honored.

`go-raml spec ...`

## Contribute

When you want to contribute to the development, follow the [contribution guidelines](contributing.md).

## roadmap
**v0.1**

* Generation of the server using [gorilla mux](http://www.gorillatoolkit.org/pkg/mux)
* Generation of a go client
* Generation of a python 3.5 client
* Generation of a python flask server

**v0.2**

* OAuth 2.0 support
* Input validation according to the RAML type definitions

**v0.3**

* Possibility to generate a default server implementation to MongoDB

**v0.4**

* Generation of a new RAML specification file

**v0.5**

* Update of a RAML specification file
