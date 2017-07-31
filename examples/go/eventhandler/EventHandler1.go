package eventhandler

import (
	"fmt"

	kl "github.com/paalth/kless/pkg/interface/klessgo"
)

//EventHandler dummy for now
type EventHandler struct {
}

//Handler the actual event handler that does nothing in this case
func (t EventHandler) Handler(c *kl.Context, resp *kl.Response, req *kl.Request) {
	fmt.Printf("Inside funcHandler\n")
}
