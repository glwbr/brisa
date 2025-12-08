package scraper

import (
	"fmt"
	"maps"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// FormState holds the ASP.NET form state fields required for POST requests.
type FormState struct {
	ViewState          string
	ViewStateGenerator string
	EventValidation    string
	// Additional ASP.NET fields that some forms require
	LastFocus     string
	EventTarget   string
	EventArgument string
}

// IsValid returns true if the form state has at least the ViewState field.
func (f *FormState) IsValid() bool {
	return f.ViewState != ""
}

// ToValues returns the form state as url.Values suitable for form submission.
func (f *FormState) ToValues() map[string]string {
	values := make(map[string]string)

	if f.ViewState != "" {
		values["__VIEWSTATE"] = f.ViewState
	}
	if f.ViewStateGenerator != "" {
		values["__VIEWSTATEGENERATOR"] = f.ViewStateGenerator
	}
	if f.EventValidation != "" {
		values["__EVENTVALIDATION"] = f.EventValidation
	}
	// Always include these fields, even if empty (some ASP.NET forms require them)
	values["__LASTFOCUS"] = f.LastFocus
	values["__EVENTTARGET"] = f.EventTarget
	values["__EVENTARGUMENT"] = f.EventArgument

	return values
}

// ParseFormState extracts ASP.NET form state from an HTML page.
func ParseFormState(html []byte) (*FormState, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return nil, fmt.Errorf("parse html: %w", err)
	}

	return ParseFormStateFromDoc(doc), nil
}

// ParseFormStateFromDoc extracts ASP.NET form state from a goquery document.
func ParseFormStateFromDoc(doc *goquery.Document) *FormState {
	getField := func(name string) string {
		val, _ := doc.Find(fmt.Sprintf("input[name='%s']", name)).Attr("value")
		return val
	}

	return &FormState{
		ViewState:          getField("__VIEWSTATE"),
		ViewStateGenerator: getField("__VIEWSTATEGENERATOR"),
		EventValidation:    getField("__EVENTVALIDATION"),
		LastFocus:          getField("__LASTFOCUS"),
		EventTarget:        getField("__EVENTTARGET"),
		EventArgument:      getField("__EVENTARGUMENT"),
	}
}

// FormBuilder helps construct form data for ASP.NET POST requests.
type FormBuilder struct {
	state  *FormState
	fields map[string]string
}

// NewFormBuilder creates a new form builder with the given form state.
func NewFormBuilder(state *FormState) *FormBuilder {
	return &FormBuilder{
		state:  state,
		fields: make(map[string]string),
	}
}

// Set adds or updates a form field.
func (b *FormBuilder) Set(key, value string) *FormBuilder {
	b.fields[key] = value
	return b
}

// SetIf adds a form field only if the condition is true.
func (b *FormBuilder) SetIf(condition bool, key, value string) *FormBuilder {
	if condition {
		b.fields[key] = value
	}
	return b
}

// Build constructs the final form values map.
func (b *FormBuilder) Build() map[string]string {
	result := make(map[string]string)

	// Add form state fields
	if b.state != nil {
		for k, v := range b.state.ToValues() {
			result[k] = v
		}
	}

	// Add custom fields
	maps.Copy(result, b.fields)

	return result
}

// FindFormAction extracts the action URL from a form element.
func FindFormAction(html []byte, formSelector string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return "", fmt.Errorf("parse html: %w", err)
	}

	action, exists := doc.Find(formSelector).Attr("action")
	if !exists {
		return "", fmt.Errorf("form action not found for selector: %s", formSelector)
	}

	return action, nil
}

// HasElement checks if an HTML page contains an element matching the selector.
func HasElement(html []byte, selector string) bool {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return false
	}
	return doc.Find(selector).Length() > 0
}

// ExtractText extracts text content from an element matching the selector.
func ExtractText(html []byte, selector string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(doc.Find(selector).Text())
}

// ExtractAttribute extracts an attribute value from an element.
func ExtractAttribute(html []byte, selector, attr string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(html)))
	if err != nil {
		return ""
	}
	val, _ := doc.Find(selector).Attr(attr)
	return val
}
