package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"time"
)

// A Handler handles log records produced by a Logger.
type Handler struct {
	attrs  []slog.Attr
	groups []string
	opts   slog.HandlerOptions
	mu     *sync.Mutex
	w      io.Writer
}

// NewHandler creates a handler that writes to w, using the given options.
// If opts is nil, the default options are used.
func NewHandler(w io.Writer, opts *slog.HandlerOptions) *Handler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &Handler{
		w:    w,
		opts: *opts,
		mu:   &sync.Mutex{},
	}
}

// Enabled reports whether the handler handles records at the given level.
// The handler ignores records whose level is lower.
func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return level >= minLevel
}

// WithAttrs returns a new handler whose attributes consists of h's attributes
// followed by attrs.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{
		attrs:  append(h.attrs, attrs...),
		groups: h.groups,
		opts:   h.opts,
		mu:     h.mu,
		w:      h.w,
	}
}

// WithGroup returns a new Handler with the given group appended to
// the receiver's existing groups.
func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{
		attrs:  h.attrs,
		groups: append(h.groups, name),
		opts:   h.opts,
		mu:     h.mu,
		w:      h.w,
	}
}

// Handle formats its argument Record so that message is followed by each
// of it's attributes on seperate lines.
func (h *Handler) Handle(_ context.Context, r slog.Record) error {
	var attrs string

	for _, a := range h.attrs {
		if !a.Equal(slog.Attr{}) {
			attrs += gray(" ↳ " + a.Key + ": " + a.Value.String() + "\n")
		}
	}

	r.Attrs(func(a slog.Attr) bool {
		if !a.Equal(slog.Attr{}) {
			attrs += gray(" ↳ " + a.Key + ": " + a.Value.String() + "\n")
		}
		return true
	})

	h.mu.Lock()
	_, err := fmt.Fprintf(h.w, "%s %s %s\n%s", r.Time.Format(time.TimeOnly), text(levelColour(r.Level), r.Level.String()), r.Message, attrs)
	h.mu.Unlock()

	return err
}

// SetDefault is syntactic sugar for constructing a new devslog handler
// and setting it as the default [slog.Logger]. The top-level slog
// functions [slog.Info], [slog.Debug], etc will all use this handler
// to format the records.
func SetDefault(w io.Writer, opts *slog.HandlerOptions) {
	slog.SetDefault(slog.New(NewHandler(w, opts)))
}
