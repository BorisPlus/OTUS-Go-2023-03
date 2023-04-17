package main_b

import (
	"testing"
	"main"

)
func TestReverse(t *testing.T) {
    got := main.Reverse("from left to right")
    if got != "thgir ot tfel morf" {
        t.Errorf("Reverse('from left to right') = %s; want 'thgir ot tfel morf'", got)
    } else {
		
	}
}
