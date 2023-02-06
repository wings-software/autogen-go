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
		drone.ConfigurePlatform,
		drone.ConfigureGo,
		drone.ConfigureNode,
		drone.ConfigurePython,
		drone.ConfigureRails,
		drone.ConfigureRuby,
		drone.ConfigureRust,
		drone.ConfigureSwift,
		drone.ConfigureDocker,
		drone.ConfigureDefault,
	}
}
