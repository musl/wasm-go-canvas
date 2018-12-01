package main

import (
	"log"
	"syscall/js"
)

const (

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/musl/wandr/lib/term"
	"github.com/nsf/termbox-go"
)

const (
	// AppName is a short name for this application.
	AppName = `wandr`
)

var (
	// Version is the semantic version string of this application. It
	// should be replaced by the build process using:
	// `-ldflags="-X main.Version="..."`
	Version = ``

	// Revision is the most recent DVCS change id. It should be replaced
	// by the build process using:
	// `-ldflags="-X main.Revision="..."`
	Revision = ``
)


func main() {
	log.SetOutput(os.Stderr)
	log.SetPrefix(AppName + ` `)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Printf("%s version: %s revision: %s", AppName, Version, Revision)


	dcb := js.NewCallback(handleDecrypt)
	defer dcb.Release()
	js.Global().Set("decrypt", dcb)

	//
	// Handle Unload
	//
	unload := make(chan struct{})
	bu := js.NewEventCallback(0, func(v js.Value) {
		unload <- struct{}{}
	})
	defer bu.Release()
	js.Global().Get("addEventListener").Invoke("beforeunload", bu)
	<-unload
}
