package builder

import (
	"github.com/wings-software/autogen-go/builder/rules"
)

type Vendor interface {
	GetRules() []rules.Rule
}

func NewVendor(vendor string) Vendor {
	switch vendor {
	case "harness":
		return rules.NewHarness()
	default:
		return rules.NewDrone()
	}
}
