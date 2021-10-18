package render

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var pageData = []struct {
	name          string
	renderer      string
	template      string
	errorExpected bool
	errorMessage  string
}{
	{"go_page", "go", "home", false, "no error expected, but error rendering go template"},
	{"go_page_no_template", "go", "no-file", true, "error expected, but no error rendering non-existent go template"},
	{"jet_page", "jet", "home", false, "no error expected, but error rendering jet template"},
	{"jet_page_no_template", "jet", "no-file", true, "error expected, but no error rendering non-existent jet template"},
	{"invalid_render_engine", "foo", "home", true, "error expected, but no error rendering with an invalid template engine"},
}

func TestRender_Page(t *testing.T) {
	for _, e := range pageData {
		r, err := http.NewRequest("GET", "/some-url", nil)
		if err != nil {
			t.Error(err)
		}

		w := httptest.NewRecorder()

		testRenderer.Renderer = e.renderer
		testRenderer.RootPath = "./testdata"

		err = testRenderer.Page(w, r, e.template, nil, nil)
		if e.errorExpected {
			if err == nil {
				t.Errorf("%s: %s", e.name, e.errorMessage)
			}
		} else {
			if err != nil {
				t.Errorf("%s: %s: %s", e.name, e.errorMessage, err.Error())
			}
		}
	}
}

func TestRender_GoPage(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/url", nil)
	if err != nil {
		t.Error(err)
	}

	testRenderer.Renderer = "go"
	testRenderer.RootPath = "./testdata"

	err = testRenderer.Page(w, r, "home", nil, nil)
	if err != nil {
		t.Error("Error rendering go template page:", err)
	}
}

func TestRender_JetPage(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/url", nil)
	if err != nil {
		t.Error(err)
	}

	testRenderer.Renderer = "jet"

	err = testRenderer.Page(w, r, "home", nil, nil)
	if err != nil {
		t.Error("Error rendering jet template page:", err)
	}
}
