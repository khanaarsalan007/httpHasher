# httpHasher

A basic tool used for creating md5 hash sum of any web based document.

## Flags

Flags are all optional, and are set with a single dash on the command line, e.g.

```
httpHasher \
-parallel  5 (default 10)

```

## Args

Website urls can be passed as args to the binary 


General usage of iterscraper:

```
httpHasher -parallel 5 http://google.com http://twitter.com
```
For an explanation of the options, type `httpHasher -help`

## URL Structure

URLs should be given in a following structure
```
http://example1.com
http://example2.com
http://example3.com
```

## Installation

Building the source requires the [Go programming language](https://golang.org/doc/install) and the [Glide](http://glide.sh) package manager.

