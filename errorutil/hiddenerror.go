package errorutil

type HiddenOriginalError struct {
	originalErr error
}

func NewHiddenOriginalError(originalErr error) *HiddenOriginalError {
	return &HiddenOriginalError{
		originalErr: originalErr,
	}
}

func (h HiddenOriginalError) Error() string {
	return ""
}

func (h HiddenOriginalError) Unwrap() error {
	return h.originalErr
}
