package scraper

import "testing"

func TestParseFormState(t *testing.T) {
	html := []byte(`
<form>
<input type="hidden" name="__VIEWSTATE" value="test-vs">
<input type="hidden" name="__VIEWSTATEGENERATOR" value="test-vsg">
<input type="hidden" name="__EVENTVALIDATION" value="test-ev">
</form>
`)

	state, err := ParseFormState(html)
	if err != nil {
		t.Fatalf("ParseFormState() error = %v", err)
	}

	if state.ViewState != "test-vs" {
		t.Errorf("ViewState = %q, want %q", state.ViewState, "test-vs")
	}
	if !state.IsValid() {
		t.Error("expected IsValid() to be true")
	}
}

func TestFormBuilder(t *testing.T) {
	state := &FormState{ViewState: "vs"}
	form := NewFormBuilder(state).
		Set("field1", "value1").
		Set("field2", "value2").
		Build()

	if form["__VIEWSTATE"] != "vs" {
		t.Errorf("__VIEWSTATE = %q, want %q", form["__VIEWSTATE"], "vs")
	}
	if form["field1"] != "value1" {
		t.Errorf("field1 = %q, want %q", form["field1"], "value1")
	}
}
