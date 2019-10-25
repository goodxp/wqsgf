package wqSGF

import (
	"testing"
)

func TestParse(t *testing.T) {
	want := "(;FF[4]C[root](;C[a];C[b](;C[c])(;C[d];C[e]))(;C[f](;C[g];C[h];C[i])(;C[j])))"

	tree := Parse(want)
	got := ToSGF(tree)

	if got != want {
		t.Errorf("want tree %v got %v", want, got)
	}
}
