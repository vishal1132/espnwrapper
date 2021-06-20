package espn

import (
	"net/http"
	"time"
)

// ESPN is the driving struct of the program.
type ESPN struct {
	// C is the http Client
	c *http.Client
}

// New returns new instance of the ESPN struct, with already configured http client.
func New() *ESPN {
	return &ESPN{
		c: &http.Client{Timeout: 30 * time.Second},
	}
}
