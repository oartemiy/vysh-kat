package novice_policy

import "vysh-kat/internal/domain"

type NovicePolicy interface {
	IsNoviceFriendly(t domain.Transport) bool
}
