package rules

import (
	"github.com/wings-software/autogen-go/builder/rules/common"
	"github.com/wings-software/autogen-go/builder/rules/harness"
)

type Harness struct{}

func NewHarness() Harness {
	return Harness{}
}

func (h Harness) GetRules() []Rule {
	return []Rule{
		harness.ConfigurePipelineVersion,
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
