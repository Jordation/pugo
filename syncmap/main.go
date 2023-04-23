package syncmap

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type UserMap struct {
	sync.RWMutex
	// strings are user ids
	ActiveLobbies map[string]uuid.UUID
}

func (um *UserMap) Get(id string) uuid.UUID {
	um.RLock()
	defer um.RUnlock()
	val, ok := um.ActiveLobbies[id]
	if ok {
		return val
	} else {
		// for my use case, should.. never happen
		panic("id not in map")
	}
}

func (u *UserMap) Set(u_uid uuid.UUID, ids ...string) {
	u.Lock()
	defer u.Unlock()
	for _, v := range ids {
		u.ActiveLobbies[v] = u_uid
	}
	fmt.Println("set these ids ", ids)
}
