package domain

import "errors"

var (
	ErrNotImplemented       = errors.New("not implemented")
	ErrInvalidCommand       = errors.New("была введена недопустимая команда")
	ErrInvalidTransportName = errors.New("неподдерживаемое имя для транспора")
	ErrInvalidThingName     = errors.New("неподдерживаемое имя для вещи")
)
