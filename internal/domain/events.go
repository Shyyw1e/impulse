package domain

import ()

type EventID int

//ingoing events
const (
    EventRegistered EventID = 1
    EventEnteredDungeon EventID = 2
    EventKilledMonster EventID = 3
    EventWentNextFloor EventID = 4
    EventWentPreviousFloor EventID = 5
    EventEnteredBossFloor EventID = 6
    EventKilledBoss EventID = 7
    EventLeftDungeon EventID = 8
    EventCannotContinue EventID = 9
    EventRestoredHealth EventID = 10
    EventReceivedDamage EventID = 11
)

//outgoing events
const (
    EventDisqualified EventID = 31
    EventDead EventID = 32
    EventImpossibleMove EventID = 33
)
