package gobootconfig

const CmdSource = "cmd"
const EnvSource = "env"

type staticPropertySource struct {
	source string
	name   string
	value  interface{}
}

func (sps staticPropertySource) getSource() string {
	return sps.source
}
func (sps staticPropertySource) getName() string {
	return sps.name
}
func (sps staticPropertySource) getValue() interface{} {
	return sps.value
}

func loadFromCmdLine() {
	cmdLine := parseCmdLine()
	for k, v := range cmdLine {
		_addPropertySource(k, staticPropertySource{CmdSource, k, v})
	}
}

func loadFromEnvironment() {
	env := parseEnv()
	for k, v := range env {
		_addPropertySource(k, staticPropertySource{EnvSource, k, v})
	}
}
