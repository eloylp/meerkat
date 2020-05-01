package unique

import (
	"log"

	"github.com/google/uuid"
)

func UUID4() string {
	uid, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
	}
	return uid.String()
}
