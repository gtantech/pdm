package models

import (
	"github.com/gtantech/pdm"
)

type Network struct {
	pdm.PDM[Activity]
}

func NewNetwork() *Network {
	return &Network{PDM: pdm.New[Activity]()}
}
