package lpac

import (
	"encoding/json"
	"fmt"
)

type Error struct {
	Message string
	Details json.RawMessage
}

func (e Error) Error() string {
	if e.Details == nil {
		return e.Message
	}
	return fmt.Sprintf("%s: %s", e.Message, string(e.Details))
}
