LilyFlux
========

LilyFlux is a small program that collects stats from a lilypad cluster.
It uses Influx DB to store the data it gathers.


Compilation
-------------

Pull the project and get the dependencies:
```bash
$ go get github.com/Tzeentchful/LilyFlux
$ go get github.com/LilyPad/GoLilyPad
$ go get github.com/influxdb/influxdb-go
$ go get launchpad.net/goyaml
```

### Then build ###

```bash
$ cd $GOPATH/pkg/github.com/Tzeentchful/LilyFlux/main
$ go build
$ go install
$ ./main
```

To run, LilyFlux will need to be able to connect to a LilyPad cluster
and have access to a Influx DB instance.
