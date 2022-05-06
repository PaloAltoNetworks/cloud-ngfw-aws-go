package permissions

import (
	"fmt"
)

// Choose returns the correct JWT style for the given scope.
func Choose(v string) (string, error) {
	switch v {
	case "", LocalScope:
		return Rulestack, nil
	case GlobalScope:
		return GlobalRulestack, nil
	}

	return "", fmt.Errorf("Unknown permission: %s", v)
}
