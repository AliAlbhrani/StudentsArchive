package helpers

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func GenerateUniqueName() string {
	return fmt.Sprintf("%s-%d", uuid.New(), time.Now().Unix())
}
