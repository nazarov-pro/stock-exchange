package time

import "testing"


func TestEpoch(t *testing.T) {
	epoch := Epoch()
	if epoch <= 0 {
		t.Fatalf("Epoch should be greater than zero. current value %d", epoch)
	}
}