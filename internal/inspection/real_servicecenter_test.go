package inspection_test

import (
	"testing"

	"vysh-kat/internal/domain"
	"vysh-kat/internal/inspection"
)

func TestRealServiceCenter_Inspect(t *testing.T) {
	tests := []struct {
		name       string
		simplicity int
		energy     float64
		accepted   bool
		reason     string
	}{
		{"валидный транспорт", 5, 1.5, true, ""},
		{"нижняя граница простоты", 1, 0, true, ""},
		{"верхняя граница простоты", 10, 100, true, ""},
		{"простота меньше единицы", 0, 1, false, "простота должна быть в диапазоне от 1 до 10"},
		{"простота больше десяти", 11, 1, false, "простота должна быть в диапазоне от 1 до 10"},
		{"отрицательная простота", -3, 1, false, "простота должна быть в диапазоне от 1 до 10"},
		{"отрицательная энергия", 5, -0.1, false, "суточный расход энергии не может быть отрицательным"},
	}

	sc := inspection.RealServiceCenter{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sc.Inspect(domain.NewScooter("Самокат", "INV-1", tt.energy, tt.simplicity))
			if got.Accepted != tt.accepted {
				t.Fatalf("accepted: ожидалось %v, получено %v (reason: %q)", tt.accepted, got.Accepted, got.Reason)
			}
			if tt.reason != "" && got.Reason != tt.reason {
				t.Fatalf("reason: ожидалось %q, получено %q", tt.reason, got.Reason)
			}
		})
	}
}
