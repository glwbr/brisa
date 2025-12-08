package scraper

import (
	"testing"
)

func TestParseFormState(t *testing.T) {
	html := []byte(`
<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
<form>
<input type="hidden" name="__VIEWSTATE" value="test-view-state">
<input type="hidden" name="__VIEWSTATEGENERATOR" value="test-generator">
<input type="hidden" name="__EVENTVALIDATION" value="test-validation">
<input type="hidden" name="__LASTFOCUS" value="">
<input type="hidden" name="__EVENTTARGET" value="">
<input type="hidden" name="__EVENTARGUMENT" value="">
</form>
</body>
</html>
`)

	state, err := ParseFormState(html)
	if err != nil {
		t.Fatalf("ParseFormState() error = %v", err)
	}

	if state.ViewState != "test-view-state" {
		t.Errorf("ViewState = %q, want %q", state.ViewState, "test-view-state")
	}
	if state.ViewStateGenerator != "test-generator" {
		t.Errorf("ViewStateGenerator = %q, want %q", state.ViewStateGenerator, "test-generator")
	}
	if state.EventValidation != "test-validation" {
		t.Errorf("EventValidation = %q, want %q", state.EventValidation, "test-validation")
	}
}

func TestFormState_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		state FormState
		want  bool
	}{
		{
			name:  "valid with all fields",
			state: FormState{ViewState: "vs", ViewStateGenerator: "vsg", EventValidation: "ev"},
			want:  true,
		},
		{
			name:  "valid with only viewstate",
			state: FormState{ViewState: "vs"},
			want:  true,
		},
		{
			name:  "invalid empty",
			state: FormState{},
			want:  false,
		},
		{
			name:  "invalid without viewstate",
			state: FormState{ViewStateGenerator: "vsg", EventValidation: "ev"},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.state.IsValid(); got != tt.want {
				t.Errorf("FormState.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormState_ToValues(t *testing.T) {
	state := FormState{
		ViewState:          "vs",
		ViewStateGenerator: "vsg",
		EventValidation:    "ev",
	}

	values := state.ToValues()

	if values["__VIEWSTATE"] != "vs" {
		t.Errorf("__VIEWSTATE = %q, want %q", values["__VIEWSTATE"], "vs")
	}
	if values["__VIEWSTATEGENERATOR"] != "vsg" {
		t.Errorf("__VIEWSTATEGENERATOR = %q, want %q", values["__VIEWSTATEGENERATOR"], "vsg")
	}
	if values["__EVENTVALIDATION"] != "ev" {
		t.Errorf("__EVENTVALIDATION = %q, want %q", values["__EVENTVALIDATION"], "ev")
	}
	// Verify additional ASP.NET fields are included (even if empty)
	if _, exists := values["__LASTFOCUS"]; !exists {
		t.Error("__LASTFOCUS should be present")
	}
	if _, exists := values["__EVENTTARGET"]; !exists {
		t.Error("__EVENTTARGET should be present")
	}
	if _, exists := values["__EVENTARGUMENT"]; !exists {
		t.Error("__EVENTARGUMENT should be present")
	}
}

func TestFormBuilder(t *testing.T) {
	state := &FormState{
		ViewState:          "vs",
		ViewStateGenerator: "vsg",
	}

	builder := NewFormBuilder(state).
		Set("field1", "value1").
		Set("field2", "value2").
		SetIf(true, "conditional", "included").
		SetIf(false, "skipped", "not-included")

	result := builder.Build()

	// Check form state fields
	if result["__VIEWSTATE"] != "vs" {
		t.Errorf("__VIEWSTATE = %q, want %q", result["__VIEWSTATE"], "vs")
	}

	// Check custom fields
	if result["field1"] != "value1" {
		t.Errorf("field1 = %q, want %q", result["field1"], "value1")
	}
	if result["field2"] != "value2" {
		t.Errorf("field2 = %q, want %q", result["field2"], "value2")
	}
	if result["conditional"] != "included" {
		t.Errorf("conditional = %q, want %q", result["conditional"], "included")
	}
	if _, exists := result["skipped"]; exists {
		t.Error("skipped field should not exist")
	}
}

func TestHasElement(t *testing.T) {
	html := []byte(`<div id="test"><span class="inner">content</span></div>`)

	tests := []struct {
		selector string
		want     bool
	}{
		{"#test", true},
		{".inner", true},
		{"#nonexistent", false},
		{".missing", false},
	}

	for _, tt := range tests {
		t.Run(tt.selector, func(t *testing.T) {
			if got := HasElement(html, tt.selector); got != tt.want {
				t.Errorf("HasElement(%q) = %v, want %v", tt.selector, got, tt.want)
			}
		})
	}
}

func TestExtractText(t *testing.T) {
	html := []byte(`<div id="test">  Hello World  </div>`)

	got := ExtractText(html, "#test")
	want := "Hello World"

	if got != want {
		t.Errorf("ExtractText() = %q, want %q", got, want)
	}
}

func TestExtractAttribute(t *testing.T) {
	html := []byte(`<a href="https://example.com" class="link">Click</a>`)

	tests := []struct {
		selector string
		attr     string
		want     string
	}{
		{"a", "href", "https://example.com"},
		{"a", "class", "link"},
		{"a", "nonexistent", ""},
	}

	for _, tt := range tests {
		t.Run(tt.attr, func(t *testing.T) {
			got := ExtractAttribute(html, tt.selector, tt.attr)
			if got != tt.want {
				t.Errorf("ExtractAttribute(%q, %q) = %q, want %q", tt.selector, tt.attr, got, tt.want)
			}
		})
	}
}
