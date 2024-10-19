package core

import "time"

type IDType int32

type Clock interface {
	Now() time.Time
}

func Core(name string) string {
	result := "Core " + name
	return result
}
