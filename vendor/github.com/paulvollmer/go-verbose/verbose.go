package verbose

import (
	"fmt"
	"io"
)

// Verbose to create simple verbose print api
type Verbose struct {
	Active bool
	Writer io.Writer
}

// New returns a Verbose object.
// the w parameter set the writer the verbose print funcs write to.
// if the a parameter was set to true, the print funcs write to the writer.
func New(w io.Writer, a bool) *Verbose {
	v := &Verbose{}
	v.Active = a
	v.Writer = w
	return v
}

// Print usage is equivalent to https://golang.org/pkg/fmt/#Print
func (v *Verbose) Print(a ...interface{}) (n int, err error) {
	if v.Active {
		n, err = fmt.Fprint(v.Writer, a...)
	}
	return
}

// Printf usage is equivalent to https://golang.org/pkg/fmt/#Printf
func (v *Verbose) Printf(format string, a ...interface{}) (n int, err error) {
	if v.Active {
		n, err = fmt.Fprintf(v.Writer, format, a...)
	}
	return
}

// Println usage is equivalent to https://golang.org/pkg/fmt/#Println
func (v *Verbose) Println(a ...interface{}) (n int, err error) {
	if v.Active {
		n, err = fmt.Fprintln(v.Writer, a...)
	}
	return
}
