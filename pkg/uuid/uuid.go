package uuid

import "github.com/google/uuid"

func NewUUID() string {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.New().String()
	}
	return id.String()
}

func NewUUIDs(n int) []string {
	ids := make([]string, n)
	for i := 0; i < n; i++ {
		ids[i] = NewUUID()
	}
	return ids
}
