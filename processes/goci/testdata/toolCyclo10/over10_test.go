package toolCyclo10

import "testing"

func TestOver10(t *testing.T) {
	result := over10()
	if !result {
		t.Errorf("over10 should be true")
	}
}
