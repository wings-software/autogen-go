package rules

import (
	"github.com/wings-software/autogen-go/builder/rules/drone"
)

type Drone struct{}

func NewDrone() Drone {
	return Drone{}
}

func (d Drone) GetRules() []Rule {
	return []Rule{
		drone.FromDrone,
	}
}
