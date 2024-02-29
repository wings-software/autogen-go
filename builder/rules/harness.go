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
		harness.ConfigurePlatform,
		harness.ConfigureGo,
		harness.ConfigureGoGlide,
		harness.ConfigureNode,
		harness.ConfigurePython,
		harness.ConfigureC,
		harness.ConfigureCPP,
		harness.ConfigureRails,
		harness.ConfigureRuby,
		harness.ConfigureRust,
		harness.ConfigureSwift,
		harness.ConfigureDefault,
		harness.ConfigureKotlin,
		harness.ConfigureKotlinwithMaven,
		harness.ConfigureDocker,
		harness.ConfigureScala,
		harness.ConfigurePHP,
	}
}
