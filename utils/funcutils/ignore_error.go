package funcutils

func IgnoreError(fn func() error) {
	_ = fn()
}
