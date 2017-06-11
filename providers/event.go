package providers

import (
	"github.com/chuckpreslar/emission"
	"github.com/olegakbarov/io.confs.core/core"
)

type (
	event struct {
		emitter *emission.Emitter
	}
)

func NewEmitter() core.Emitter {
	return &event{emission.NewEmitter()}
}

func (e *event) On(event, listener interface{}) {
	e.emitter.On(event, listener)
}

func (e *event) Emit(event interface{}, args ...interface{}) {
	e.emitter.Emit(event, args...)
}
