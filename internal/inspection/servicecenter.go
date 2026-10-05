package inspection

import "vysh-kat/internal/domain"

type ResultOfInspection struct {
	Accepted bool
	Reason   string
}

type ServiceCenter interface {
	Inspect(t domain.Transport) ResultOfInspection
}
