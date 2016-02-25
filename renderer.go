package static

type Renderer interface {
	Start(s Static, ev EventHandler) error
}
