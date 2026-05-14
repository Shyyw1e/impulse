package domain

import (
	"time"

	"github.com/Shyyw1e/impulse/internal/config"
)

type Dungeon struct {
	FloorsWithMonsters 		int
	MonstersPerFloor		int
	MonstersKilledByFloor	[]int
	BossKilled				bool
}

func NewDungeon(cfg *config.Config) *Dungeon {
	FloorsWithMonsters:= cfg.Floors - 1

	if FloorsWithMonsters == 0 {
		return &Dungeon{BossKilled: false}
	}

	dunge := Dungeon{
		FloorsWithMonsters: cfg.Floors - 1,
		MonstersPerFloor: cfg.Monsters,
		BossKilled: false,
	}
	
	killedMonsters := make([]int, dunge.FloorsWithMonsters)
	for i := 0; i < dunge.MonstersPerFloor; i++ {
		killedMonsters[i] = dunge.MonstersPerFloor
	}

	dunge.MonstersKilledByFloor = killedMonsters

	return &dunge
}

type Player struct {
	ID		int
	Registered 	bool
	InDungeon	bool
	Floor		int
	HP			int
	Status		TrialStatus

	EnteredAt   time.Time
	FinishedAt	time.Time
}