package database

import "testing"

func TestOpenRejectsMalformedURL(t *testing.T) {
	if _, err := Open("postgres://%zz"); err == nil {
		t.Fatal("expected malformed connection URL to return an error")
	}
}
