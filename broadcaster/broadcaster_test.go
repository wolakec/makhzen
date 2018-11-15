package broadcaster

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendMessage(t *testing.T) {
	b := Broadcaster{}

	t.Run("SendMessage recieves an error from unreachable host", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusServiceUnavailable)
		}))
		defer ts.Close()

		url := ts.URL

		err := b.SendMessage("region", "us-east-1", url)

		if err == nil {
			t.Error("SendMessage did not return an error")
		}
	})

	t.Run("SendMessage makes valid POST request", func(t *testing.T) {
		wantKey := "region"
		wantValue := "eu-west-1"

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)

			if r.Method != "POST" {
				t.Errorf("expected method POST, got %s", r.Method)
			}

			expectedPath := "/message"
			if r.URL.EscapedPath() != expectedPath {
				t.Errorf("expected path %s, got %s", expectedPath, r.URL.EscapedPath())
			}

			r.ParseForm()

			b, err := ioutil.ReadAll(r.Body)

			if err != nil {
				t.Errorf("could not read request body: %s", err)
			}

			var msg Message
			err = json.Unmarshal(b, &msg)

			if err != nil {
				t.Errorf("could not parse request body into message %s", err)
			}

			if msg.Key != wantKey {
				t.Errorf("incorrect key, expected: %s, got: %s", wantKey, msg.Key)
			}

			if msg.Value != wantValue {
				t.Errorf("incorrect value, expected: %s, got: %s", wantValue, msg.Value)
			}
		}))
		defer ts.Close()

		url := ts.URL

		err := b.SendMessage(wantKey, wantValue, url)

		if err != nil {
			t.Errorf("SendMessage returned error: %s", err)
		}
	})
}
