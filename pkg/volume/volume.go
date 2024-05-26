package main

import (
	"fmt"
	"log"
	"net"

	"github.com/go-git/go-billy/v5/memfs"
	nfs "github.com/willscott/go-nfs"
	nfshelper "github.com/willscott/go-nfs/helpers"
)

func main() {
	listener, err := net.Listen("tcp", ":0")
	panicOnErr(err, "starting TCP listener")
	fmt.Printf("Server running at %s\n", listener.Addr())
	mem := memfs.New()
	f, err := mem.Create("hello.txt")
	panicOnErr(err, "creating file")
	_, err = f.Write([]byte("hello world"))
	panicOnErr(err, "writing data")
	f.Close()
	handler := nfshelper.NewNullAuthHandler(mem)
	cacheHelper := nfshelper.NewCachingHandler(handler, 1)
	panicOnErr(nfs.Serve(listener, cacheHelper), "serving nfs")
}

func panicOnErr(err error, desc ...interface{}) {
	if err == nil {
		return
	}
	log.Println(desc...)
	log.Panicln(err)
}
