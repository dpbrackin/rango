package core

import "github.com/google/uuid"

type IDType uuid.UUID

// MarshalText implements encoding.TextMarshaler.
func (id IDType) MarshalText() ([]byte, error) {
	return uuid.UUID(id).MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (id *IDType) UnmarshalText(data []byte) error {
	i, err := uuid.ParseBytes(data)
	if err != nil {
		return err
	}
	*id = IDType(i)
	return nil
}

func (id *IDType) String() string {
	return uuid.UUID(*id).String()
}
