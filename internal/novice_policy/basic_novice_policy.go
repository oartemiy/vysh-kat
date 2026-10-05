package novice_policy

import "vysh-kat/internal/domain"

type BasicNovicePolicy struct{}

func (p *BasicNovicePolicy) IsNoviceFriendly(t domain.Transport) bool {
	return t.GetSimplicity() >= 6
}
