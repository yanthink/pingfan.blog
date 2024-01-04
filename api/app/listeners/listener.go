package listeners

type listener struct {
}

func (*listener) ShouldSync() bool {
	return false
}
