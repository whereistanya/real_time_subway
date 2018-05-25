# Real time subway info

## install go
$ sudo apt-get install golang  
$ mkdir $HOME/code/go  
$ export GOPATH=$HOME/code/go  

$ sudo apt-get install protobuf-compiler

$ go get -u -v github.com/golang/protobuf/proto  
$ go get -u -v github.com/golang/protobuf/protoc-gen-go  

Add $GOPATH/bin to your path (because you need to access protoc-gen-go from
 there)

## Configure the key, trains and stations.
Generate a key at http://datamine.mta.info/ . Create a file called cmd/keys.go
that looks like

```
package main

// Generated a key at http://datamine.mta.info/
const KEY = "insert key here"
```

Add the feed, train and station to the top of main.go


# Protos
These come from:
https://developers.google.com/transit/gtfs-realtime/gtfs-realtime-proto
http://datamine.mta.info/sites/all/files/pdfs/nyct-subway.proto.txt

I built them with:

$ protoc --go_out=. gtfs-realtime.proto  
$ protoc --go_out=. nyct-subway.proto  

This creates gtfs-realtime.pb.go and nyct-subway.pb.go.

