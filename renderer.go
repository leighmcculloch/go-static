package static

type Renderer interface {
	Render(s Static, ev EventHandler) error
}
