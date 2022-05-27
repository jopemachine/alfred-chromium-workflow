package src

import (
	"github.com/deanishe/awgo"
)

type WorkflowConfig struct {
	Browser string
	Locale string
	Profile string
	ExcludeDomains []string
	ResultLimitCount uint8
}

var wf *aw.Workflow

var Conf = &WorkflowConfig{}

var ImportConfig = func () {
	conf := aw.NewConfig()

	if err := conf.To(Conf); err != nil {
		panic(err)
	}
}

