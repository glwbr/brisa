package scraper

import (
	"fmt"
	"maps"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// FormState holds ASP.NET form state fields.
type FormState struct {
	ViewState          string
	ViewStateGenerator string
	EventValidation    string
	LastFocus          string
	EventTarget        string
	EventArgument      string
}

func (f *FormState) IsValid() bool { return f.ViewState != "" }

func (f *FormState) Values() map[string]string {
	m := make(map[string]string)
	if f.ViewState != "" {
		m["__VIEWSTATE"] = f.ViewState
	}
	if f.ViewStateGenerator != "" {
		m["__VIEWSTATEGENERATOR"] = f.ViewStateGenerator
	}
	if f.EventValidation != "" {
		m["__EVENTVALIDATION"] = f.EventValidation
	}
	m["__LASTFOCUS"] = f.LastFocus
	m["__EVENTTARGET"] = f.EventTarget
	m["__EVENTARGUMENT"] = f.EventArgument
	return m
}

// ParseFormState extracts ASP.NET form state from HTML.
func ParseFormState(html []byte) (*FormState, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return nil, err
	}
	get := func(name string) string {
		val, _ := doc.Find(fmt.Sprintf("input[name='%s']", name)).Attr("value")
		return val
	}
	return &FormState{
		ViewState:          get("__VIEWSTATE"),
		ViewStateGenerator: get("__VIEWSTATEGENERATOR"),
		EventValidation:    get("__EVENTVALIDATION"),
		LastFocus:          get("__LASTFOCUS"),
		EventTarget:        get("__EVENTTARGET"),
		EventArgument:      get("__EVENTARGUMENT"),
	}, nil
}

// FormBuilder constructs form data for POST requests.
type FormBuilder struct {
	state  *FormState
	fields map[string]string
}

func NewFormBuilder(state *FormState) *FormBuilder {
	return &FormBuilder{state: state, fields: make(map[string]string)}
}

func (b *FormBuilder) Set(key, value string) *FormBuilder {
	b.fields[key] = value
	return b
}

func (b *FormBuilder) Build() map[string]string {
	result := make(map[string]string)
	if b.state != nil {
		maps.Copy(result, b.state.Values())
	}
	maps.Copy(result, b.fields)
	return result
}
