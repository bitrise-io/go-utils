package timeout

import "time"

// HandlerModel ....
type HandlerModel struct {
	timer     *time.Timer
	timeout   time.Duration
	onTimeout func()

	running bool
}

// Start ...
func (handler *HandlerModel) Start() {
	handler.timer = time.NewTimer(handler.timeout)
	handler.running = true

	go func() {
		for _ = range handler.timer.C {
			if handler.onTimeout != nil {
				handler.onTimeout()
			}
		}
	}()
}

// Stop ...
func (handler *HandlerModel) Stop() bool {
	if handler.running {
		handler.running = false
		return handler.timer.Stop()
	}
	return false
}

// Running ...
func (handler HandlerModel) Running() bool {
	return handler.running
}

// NewTimeoutHandler ...
func NewTimeoutHandler(timeout time.Duration, onTimeout func()) HandlerModel {
	return HandlerModel{
		timeout:   timeout,
		onTimeout: onTimeout,
	}
}
