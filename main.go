package main

import (
	"log"
	"math"
	"os"
	"strings"
	"syscall/js"
	"time"

	"github.com/alecthomas/template"
)

const (
	// AppName is a short string that identifies this application.
	AppName = `wasm-go-canvas`
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

func draw(args []js.Value) {
	if len(args) != 2 {
		log.Printf("draw needs 2 arguments, an id to draw to, and a number of frames.")
		return
	}

	doc := js.Global().Get("document")
	id := args[0].String()
	frames := args[1].Int()
	svg := doc.Call("getElementById", id)

	// 1000 / 60 = 16.666...
	step := math.Pi / 120.0
	frameDelay := 16666 * time.Nanosecond
	tmplString := `<circle cx="{{.X}}" cy="{{.Y}}" r="25" stroke="#ddd" stroke-width="1" fill="#def" />`
	tmplData := struct {
		X int
		Y int
	}{}
	tmpl, err := template.New("svg").Parse(tmplString)
	if err != nil {
		log.Printf("Unable to compile the SVG template: %s", err.Error())
		return
	}
	var tmplWriter strings.Builder

	for i := 0; i < frames; i++ {
		tmplData.X = int(250 + 200*math.Cos(float64(i)*step))
		tmplData.Y = int(250 + 200*math.Sin(float64(i)*step))

		err = tmpl.Execute(&tmplWriter, tmplData)
		if err != nil {
			log.Printf("Unable to render the SVG template: %s", err.Error())
			return
		}

		svg.Set("innerHTML", tmplWriter.String())
		tmplWriter.Reset()

		time.Sleep(frameDelay)
	}
}

func main() {
	log.SetOutput(os.Stderr)
	log.SetPrefix(AppName + ` `)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Printf("%s version: %s revision: %s", AppName, Version, Revision)

	//
	// Setup Callbacks
	//
	cb := js.NewCallback(draw)
	defer cb.Release()
	js.Global().Set("draw", cb)

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
