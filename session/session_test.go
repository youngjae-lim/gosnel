package session

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/alexedwards/scs/v2"
)

func TestSession_InitSession(t *testing.T) {
	c := &Session{
		CookieName:     "gosnel",
		CookieLifetime: "100",
		CookiePersist:  "true",
		CookieDomain:   "localhost",
		CookieSecure:   "false",
		SessionType:    "cookie",
	}

	var sm *scs.SessionManager

	sess := c.InitSession()

	var sessKind reflect.Kind
	var sessType reflect.Type

	rv := reflect.ValueOf(sess)

	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		fmt.Println("For loop:", rv.Kind(), rv.Type(), rv)
		sessKind = rv.Kind()
		sessType = rv.Type()

		rv = rv.Elem()
	}

	if !rv.IsValid() {
		t.Error("invalid type or kind; kind:", rv.Kind(), "type:", rv.Type())
	}

	if sessKind != reflect.ValueOf(sm).Kind() {
		t.Error("wrong kind returned. Expected:", reflect.ValueOf(sm).Kind(), "but got", sessKind)
	}

	if sessType != reflect.ValueOf(sm).Type() {
		t.Error("wrong type returned. Expected:", reflect.ValueOf(sm).Type(), "but got", sessType)
	}
}
