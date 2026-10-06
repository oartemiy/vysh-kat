package inspection

import "vysh-kat/internal/domain"

type RealServiceCenter struct{}

func (s RealServiceCenter) Inspect(t domain.Transport) ResultOfInspection {
	if t.GetSimplicity() <= 0 || t.GetSimplicity() >= 11 {
		return ResultOfInspection{false, "простота должна быть в диапазоне от 1 до 10"}
	}
	if t.DailyEnergyKWh() < 0 {
		return ResultOfInspection{false, "суточный расход энергии не может быть отрицательным"}
	}
	return ResultOfInspection{true, "ok"}
}
