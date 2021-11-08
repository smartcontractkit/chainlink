package gracefulpanic

func WrapRecover(fn func()) {
	// Hotfix: removed recover to ease debugging
	fn()
}
