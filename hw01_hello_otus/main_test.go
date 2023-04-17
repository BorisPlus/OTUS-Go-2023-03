package main

import (
	"testing"
)

func TestItEcho(t *testing.T) {
	got := Echo("a")
	if got != "a" {
		t.Errorf("Echo('a') = %s; want 'a'", got)
	} else {

	}
}
