package progress

// Wrapper ...
type Wrapper struct {
	spinner        Spinner
	action         func()
	failableAction func() error
}

// NewWrapper ...
func NewWrapper(spinner Spinner) Wrapper {
	return Wrapper{
		spinner: spinner,
	}
}

// NewDefaultWrapper ...
func NewDefaultWrapper(message string) Wrapper {
	return NewWrapper(NewDefaultSpinner(message))
}

// WrapAction ...
func (w Wrapper) WrapAction(action func()) {
	w.spinner.Start()
	action()
	w.spinner.Stop()
}
