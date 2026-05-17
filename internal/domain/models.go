package domain

import (
	"time"

	"github.com/Shyyw1e/impulse/internal/config"
)

type PlayerState struct {
	ID         int
	Registered bool
	InDungeon  bool
	Floor      int
	HP         int
	Status     TrialStatus
	Finished   bool

	MonstersKilledByFloor []int
	BossKilled            bool

	EnteredAt  time.Time
	FinishedAt time.Time

	FloorEnterAt        time.Time
	FloorClearDurations []time.Duration

	BossEnterAt      time.Time
	BossKillDuration time.Duration
}

func NewPlayer(id int, cfg *config.Config) *PlayerState {
	return &PlayerState{
		ID:                    id,
		HP:                    100,
		Floor:                 0,
		Status:                StatusFail,
		MonstersKilledByFloor: make([]int, cfg.Floors-1),
	}
}
