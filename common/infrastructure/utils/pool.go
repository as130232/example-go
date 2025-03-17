package utils

import "github.com/panjf2000/ants/v2"

func NewPool(size int) *ants.Pool {
	pool, err := ants.NewPool(size)
	if err != nil {
		panic(err)
	}

	return pool
}
