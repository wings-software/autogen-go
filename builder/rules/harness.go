package rules

type Harness struct{}

func NewHarness() Drone {
	return Drone{}
}

func (h Harness) GetRules() []Rule {
	return []Rule{
		// FromDrone,
		// ConfigurePlatform,
		// ConfigureGo,
		// ConfigureNode,
		// ConfigurePython,
		// ConfigureRails,
		// ConfigureRuby,
		// ConfigureRust,
		// ConfigureSwift,
		// ConfigureDocker,
		// ConfigureDefault,
	}
}
