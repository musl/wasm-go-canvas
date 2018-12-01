package main

import (
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"syscall/js"
	"time"

	template "github.com/alecthomas/template"
	colorful "github.com/lucasb-eyer/go-colorful"
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

// Element is
type Element interface {
	Update()
	String() string
}

// Context is everything we need to animate our SVG canvas.
type Context struct {
	Canvas   js.Value
	Document js.Value
	Elements []Element
	FPS      float64
}

// NewContext returns a new Context struct
func NewContext(fps float64, selector string) *Context {
	doc := js.Global().Get(`document`)
	c := &Context{
		FPS:      fps,
		Document: doc,
		Canvas:   doc.Call(`querySelector`, selector),
	}

	return c
}

func setup(c *Context) {
}

// Spinner is.
type Spinner struct {
	A           float64
	Fill        string
	I           float64
	IMax        float64
	O           float64
	Opacity     float64
	P           float64
	R           float64
	Step        float64
	Stroke      string
	StrokeWidth float64
	T           float64
	X           int
	Y           int
}

// Update moves a Spinner.
func (s *Spinner) Update() {
	s.T += s.Step
	s.P = ((4.0 * math.Pi / s.IMax) * s.I) + math.Sin((s.I/s.IMax)*s.T)
	s.X = int(s.O + s.A*math.Cos(s.T+s.P))
	s.Y = int(s.O + s.A*math.Sin(s.T+s.P))
}

func (s *Spinner) String() string {
	tmplString := `<circle cx="{{.X}}" cy="{{.Y}}" r="{{.R}}" opacity="{{.Opacity}}" stroke="{{.Stroke}}" stroke-width="{{.StrokeWidth}}" fill="{{.Fill}}" />`
	tmpl, err := template.New("svg").Parse(tmplString)
	if err != nil {
		log.Printf("Unable to compile the SVG template: %s", err.Error())
		return ""
	}
	var tmplWriter strings.Builder

	err = tmpl.Execute(&tmplWriter, s)
	if err != nil {
		log.Printf("Unable to render the SVG template: %s", err.Error())
		return ""
	}
	return tmplWriter.String()
}

func loop(c *Context) {
	delay := time.Duration(float64(time.Second) / c.FPS)
	for {
		svg := make([]string, len(c.Elements))
		for i := range c.Elements {
			c.Elements[i].Update()
			svg[i] = c.Elements[i].String()
		}
		c.Canvas.Set("innerHTML", strings.Join(svg, ``))
		time.Sleep(delay)
	}
}

func main() {
	log.SetOutput(os.Stderr)
	log.SetPrefix(AppName + ` `)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Printf("%s version: %s revision: %s", AppName, Version, Revision)

	c := NewContext(60, `svg#canvas`)
	w, err := strconv.Atoi(c.Canvas.Call(`getAttribute`, `width`).String())
	if err != nil {
		return
	}

	iMax := 16
	for i := 0; i < iMax; i++ {
		c.Elements = append(c.Elements,
			&Spinner{
				I:           float64(i),
				IMax:        float64(iMax),
				A:           10 + (float64(w)/3)*(float64(i)/float64(iMax)),
				O:           (float64(w) / 2.0),
				R:           16.0 + (float64(w)/32.0)*(float64(i)/float64(iMax)),
				Opacity:     0.5,
				Step:        math.Pi / 60.0,
				StrokeWidth: 1.0,
				Stroke:      colorful.Hsv(360*(float64(i)/float64(iMax)), 1.0, 1.0).Hex(),
				Fill:        colorful.Hsv(360*(float64(i)/float64(iMax)), 0.2, 1.0).Hex(),
			},
		)
	}
	go loop(c)

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
