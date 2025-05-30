package entities

import (
	"time"

	"github.com/google/uuid"
)

type Blacklist struct {
	ID 			uuid.UUID `json:"id"`
	IP  		string    `json:"ip"`
	Name 		string    `json:"name"`
	LastSeen 	time.Time `json:"last_seen"`
	Dispositive	string    `json:"dispositive"`
}
