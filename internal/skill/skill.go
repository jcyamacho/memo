package skill

import _ "embed"

//go:embed SKILL.md
var guide string

func Guide() string {
	return guide
}
