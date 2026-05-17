package domain

type TrialStatus string
type DungeonStatus string

const(
	StatusSuccess 	TrialStatus = "SUCCESS"
	StatusFail 		TrialStatus = "FAIL"
	StatusDisqual 	TrialStatus = "DISQUAL"
)

const (
	StatusOpen		DungeonStatus = "OPEN"
	StatusClosed	DungeonStatus = "CLOSED"
)