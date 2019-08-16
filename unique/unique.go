package unique

import (
	"github.com/google/uuid"
	"log"
)

func UUID4() string {
	uid, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
	}
	return uid.String()
}
