package novice_policy_test

import (
	"testing"

	"vysh-kat/internal/domain"
	"vysh-kat/internal/novice_policy"
)

func TestBasicNovicePolicy_IsNoviceFriendly(t *testing.T) {
	tests := []struct {
		simplicity int
		want       bool
	}{
		{1, false},
		{5, false},
		{6, true},
		{7, true},
		{10, true},
	}

	policy := &novice_policy.BasicNovicePolicy{}
	for _, tt := range tests {
		got := policy.IsNoviceFriendly(domain.NewScooter("Самокат", "INV-1", 1, tt.simplicity))
		if got != tt.want {
			t.Fatalf("simplicity=%d: ожидалось %v, получено %v", tt.simplicity, tt.want, got)
		}
	}
}
