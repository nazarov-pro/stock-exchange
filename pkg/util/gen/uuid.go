package gen

import 	"github.com/google/uuid"

// NewUUID generates new uuid
func NewUUID() string {
	return uuid.New().String()
}