package inspection

import "vysh-kat/internal/domain"

type RealServiceCenter struct{}

func (s RealServiceCenter) Inspect(t domain.Transport) ResultOfInspection {
	if t.GetSimplicity() <= 0 || t.GetSimplicity() >= 11 {
		return ResultOfInspection{false, "invalid simplicity number"}
	}
	if t.DailyEnergyKWh() < 0 {
		return ResultOfInspection{false, "Kwh is negative"}
	}
	return ResultOfInspection{true, "ok"}
}
