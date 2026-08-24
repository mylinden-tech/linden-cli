// Package output handles all CLI output formatting.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/itchyny/gojq"
)

// Format specifies the output format.
type Format int

const (
	FormatAuto Format = iota // TTY → Styled, non-TTY → JSON
	FormatJSON
	FormatMarkdown
	FormatStyled
	FormatQuiet
)

// Response is the success JSON envelope.
type Response struct {
	OK          bool           `json:"ok"`
	Data        any            `json:"data,omitempty"`
	Summary     string         `json:"summary,omitempty"`
	Notice      string         `json:"notice,omitempty"`
	Breadcrumbs []Breadcrumb   `json:"breadcrumbs,omitempty"`
	Meta        map[string]any `json:"meta,omitempty"`

	// Not serialized — used by presenter for styled/markdown output.
	DisplayData any `json:"-"`
}

// ErrorResponse is the failure JSON envelope.
type ErrorResponse struct {
	OK    bool           `json:"ok"`
	Error string         `json:"error"`
	Code  string         `json:"code"`
	Hint  string         `json:"hint,omitempty"`
	Meta  map[string]any `json:"meta,omitempty"`
}

// Breadcrumb is a suggested follow-up command.
type Breadcrumb struct {
	Action      string `json:"action"`
	Cmd         string `json:"cmd"`
	Description string `json:"description"`
}

// Options controls output behaviour.
type Options struct {
	Format    Format
	Writer    io.Writer
	ErrWriter io.Writer
	JQFilter  string
}

// ResponseOption configures a Response.
type ResponseOption func(*Response)

// WithSummary sets the summary field.
func WithSummary(s string) ResponseOption { return func(r *Response) { r.Summary = s } }

// WithNotice sets the notice field.
func WithNotice(s string) ResponseOption { return func(r *Response) { r.Notice = s } }

// WithDisplayData sets alternate data for styled/markdown rendering.
func WithDisplayData(d any) ResponseOption { return func(r *Response) { r.DisplayData = d } }

// WithBreadcrumbs appends breadcrumbs to the response.
func WithBreadcrumbs(bc ...Breadcrumb) ResponseOption {
	return func(r *Response) { r.Breadcrumbs = append(r.Breadcrumbs, bc...) }
}

// Writer handles all output formatting.
type Writer struct {
	opts Options
	jq   *gojq.Code
}

// New creates a new output Writer. Returns an error if JQFilter is set but invalid.
func New(opts Options) (*Writer, error) {
	if opts.Writer == nil {
		opts.Writer = os.Stdout
	}
	if opts.ErrWriter == nil {
		opts.ErrWriter = os.Stderr
	}

	w := &Writer{opts: opts}

	if opts.JQFilter != "" {
		q, err := gojq.Parse(opts.JQFilter)
		if err != nil {
			return nil, ErrJQValidation(err)
		}
		code, err := gojq.Compile(q)
		if err != nil {
			return nil, ErrJQValidation(err)
		}
		w.jq = code
	}
	return w, nil
}

// DefaultOptions returns options for standard output.
func DefaultOptions() Options {
	return Options{
		Format:    FormatAuto,
		Writer:    os.Stdout,
		ErrWriter: os.Stderr,
	}
}

// OK writes a success response.
func (w *Writer) OK(data any, opts ...ResponseOption) error {
	resp := &Response{OK: true, Data: data}
	for _, o := range opts {
		o(resp)
	}
	return w.write(resp)
}

// Err writes an error response envelope.
func (w *Writer) Err(err error) error {
	e := AsError(err)
	resolved := w.resolveFormat()

	if resolved == FormatStyled || (resolved == FormatAuto && isTTY(w.opts.Writer)) {
		fmt.Fprintf(w.opts.ErrWriter, "Error: %s\n", e.Message)
		if e.Hint != "" {
			fmt.Fprintf(w.opts.ErrWriter, "Hint:  %s\n", e.Hint)
		}
		return err
	}

	resp := ErrorResponse{
		OK:    false,
		Error: e.Message,
		Code:  e.Code,
		Hint:  e.Hint,
	}
	enc := json.NewEncoder(w.opts.Writer)
	enc.SetIndent("", "  ")
	_ = enc.Encode(resp)
	return err
}

func (w *Writer) write(resp *Response) error {
	if w.jq != nil {
		return w.writeJQ(resp)
	}

	switch w.resolveFormat() {
	case FormatQuiet:
		return encodeJSON(w.opts.Writer, resp.Data)
	case FormatMarkdown:
		return w.writeMarkdown(resp)
	case FormatStyled:
		return w.writeStyled(resp)
	default:
		return encodeJSON(w.opts.Writer, resp)
	}
}

func (w *Writer) resolveFormat() Format {
	if w.opts.Format == FormatAuto {
		if isTTY(w.opts.Writer) {
			return FormatStyled
		}
		return FormatJSON
	}
	return w.opts.Format
}

func (w *Writer) writeJQ(resp *Response) error {
	raw, err := jsonRoundtrip(resp.Data)
	if err != nil {
		return err
	}
	iter := w.jq.Run(raw)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if runErr, ok := v.(error); ok {
			return runErr
		}
		b, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return err
		}
		fmt.Fprintln(w.opts.Writer, string(b))
	}
	return nil
}

func (w *Writer) writeMarkdown(resp *Response) error {
	if resp.Summary != "" {
		fmt.Fprintf(w.opts.Writer, "## %s\n\n", resp.Summary)
	}
	return encodeJSON(w.opts.Writer, resp.Data)
}

func (w *Writer) writeStyled(resp *Response) error {
	if resp.Summary != "" {
		fmt.Fprintf(w.opts.Writer, "%s\n", resp.Summary)
	}
	return encodeJSON(w.opts.Writer, resp.Data)
}

func encodeJSON(w io.Writer, data any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

func jsonRoundtrip(v any) (any, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var out any
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func isTTY(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
