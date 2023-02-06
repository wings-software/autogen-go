package rules

import (
	"github.com/wings-software/autogen-go/builder/rules/common"
	"github.com/wings-software/autogen-go/builder/rules/drone"
)

type Drone struct{}

func NewDrone() Drone {
	return Drone{}
}

func (d Drone) GetRules() []Rule {
	return []Rule{
		drone.FromDrone,
		common.ConfigurePlatform,
		common.ConfigureGo,
		common.ConfigureNode,
		common.ConfigurePython,
		common.ConfigureRails,
		common.ConfigureRuby,
		common.ConfigureRust,
		common.ConfigureSwift,
		common.ConfigureDocker,
		common.ConfigureDefault,
	}
}
