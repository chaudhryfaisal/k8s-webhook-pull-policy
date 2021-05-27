package webhook

import (
	"net/http"
)

// routes wires the routes to handlers on a specific router.
func (h handler) routes(router *http.ServeMux) error {
	mark, err := h.mark()
	if err != nil {
		return err
	}
	router.Handle("/wh/mutating/mark", mark)

	return nil
}
