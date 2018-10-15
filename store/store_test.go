package store

import "testing"

func TestSet(t *testing.T) {
	var s = New()

	got := s.Set("some-key", "1234")
	want := "1234"

	if got != want {
		t.Errorf("Set was incorrect, expected %s but got %s", want, got)
	}
}

func TestSetUpdate(t *testing.T) {
	var s = New()
	s.Set("some-key", "1234")

	got := s.Set("some-key", "new-val")
	want := "new-val"

	if got != want {
		t.Errorf("Set was incorrect, expected %s but got %s", want, got)
	}
}

func TestGet(t *testing.T) {
	var s = New()
	s.Set("some-key", "1234")

	got, _ := s.Get("some-key")
	want := "1234"

	if got != want {
		t.Errorf("Get was incorrect, expected %s but got %s", want, got)
	}
}

func TestGetOnEmptyKey(t *testing.T) {
	var s = New()
	_, got := s.Get("some-key")
	want := false

	if got != want {
		t.Errorf("Get was incorrect, expected %v but got %v", want, got)
	}
}
