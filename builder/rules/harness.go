package rules

import (
	"github.com/wings-software/autogen-go/builder/rules/harness"
)

type Harness struct{}

func NewHarness() Harness {
	return Harness{}
}

func (h Harness) GetRules() []Rule {
	return []Rule{
		harness.ConfigurePipelineVersion,
		harness.ConfigurePlatform,
		harness.ConfigureGo,
		harness.ConfigureNode,
		harness.ConfigurePython,
		harness.ConfigureRails,
		harness.ConfigureRuby,
		harness.ConfigureRust,
		harness.ConfigureSwift,
		harness.ConfigureDocker,
		harness.ConfigureDefault,
		harness.ConfigureScala,
	}
}
