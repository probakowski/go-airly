# go-airly
go-airly is a Go client library for accessing the [Airly API](https://airly.org/en/pricing/airly-api/)

[![Build](https://github.com/probakowski/go-airly/actions/workflows/build.yml/badge.svg)](https://github.com/probakowski/go-airly/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/probakowski/go-airly)](https://goreportcard.com/report/github.com/probakowski/go-airly)

## Installation
go-airly is compatible with modern Go releases in module mode, with Go installed:

```bash
go get github.com/probakowski/go-airly
```

will resolve and add the package to the current development module, along with its dependencies.

Alternatively the same can be achieved if you use import in a package:

```go
import "github.com/probakowski/go-airly"
```

and run `go get` without parameters.

Finally, to use the top-of-trunk version of this repo, use the following command:

```bash
go get github.com/probakowski/go-airly@master
```

## Usage ##

```go
import "github.com/probakowski/go-airly"
```

Construct a new Airly client, then you can use different method from [API](https://airly.org/en/pricing/airly-api/), for example:

```go
client := airly.Client{
Key:        "<your API key>", //required
Language:   "pl",             //optional, options: en, pl, default en
HttpClient: client,           //optional, HTTP client to use, http.DefaultClient will be used if nil
}
installations, err := client.NearestInstallations(airly.Location{Latitude: 50.062006, Longitude: 19.940984})
...
```