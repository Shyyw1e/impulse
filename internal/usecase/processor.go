package usecase

import (
	"fmt"
	"strconv"

	"github.com/Shyyw1e/impulse/internal/config"
	"github.com/Shyyw1e/impulse/internal/domain"
)

type Processor struct {
	cfg     *config.Config
	players map[int]*domain.PlayerState
	output  []config.Event
}

func NewProcessor(cfg *config.Config) *Processor {
	return &Processor{
		cfg:     cfg,
		players: make(map[int]*domain.PlayerState),
	}
}

func (p *Processor) Process(events []config.Event) []config.Event {
	for _, ev := range events {
		p.processEvent(ev)
	}

	return p.output
}

func (p *Processor) processEvent(ev config.Event) {
	switch ev.ID {
	case config.EventRegistered:
		p.handleRegister(ev)
	case config.EventEnteredDungeon:
		p.handleEnterDungeon(ev)
	case config.EventKilledMonster:
		p.handleKilledMonster(ev)
	case config.EventWentNextFloor:
		p.handleWentNextFloor(ev)
	case config.EventWentPreviousFloor:
		p.handleWentPreviousFloor(ev)
	case config.EventEnteredBossFloor:
		p.handleEnteredBossFloor(ev)
	case config.EventKilledBoss:
		p.handleKilledBoss(ev)
	case config.EventLeftDungeon:
		p.handleLeftDungeon(ev)
	case config.EventCannotContinue:
		p.handleCannotContinue(ev)
	case config.EventRestoredHealth:
		p.handleRestoredHealth(ev)
	case config.EventReceivedDamage:
		p.handleReceivedDamage(ev)
	default:
		p.handleImpossibleMove(ev)
	}
}

func (p *Processor) getOrCreatePlayer(playerID int) *domain.PlayerState {
	if v, ok := p.players[playerID]; ok {
		return v
	}

	newPlayer := domain.NewPlayer(playerID, p.cfg)
	p.players[playerID] = newPlayer
	return newPlayer
}

func (p *Processor) canHandleDungeonEvent(ev config.Event, player *domain.PlayerState) bool {
	if !player.Registered {
		p.handleDisqual(ev)
		return false
	}
	if player.Finished || !player.InDungeon {
		p.handleImpossibleMove(ev)
		return false
	}
	if p.isDungeonClosed(ev) {
		player.Finished = true
		player.InDungeon = false
		player.FinishedAt = p.cfg.OpenAt.Add(p.cfg.Duration)
		player.Status = domain.StatusFail
		p.handleImpossibleMove(ev)
		return false
	}

	return true
}

func (p *Processor) isDungeonClosed(ev config.Event) bool {
	return !ev.Time.Before(p.cfg.OpenAt.Add(p.cfg.Duration))
}

func (p *Processor) isMonsterFloor(floor int) bool {
	return floor > 0 && floor < p.cfg.Floors
}

func (p *Processor) isCurrentFloorCleared(player *domain.PlayerState) bool {
	if !p.isMonsterFloor(player.Floor) {
		return false
	}

	return player.MonstersKilledByFloor[player.Floor-1] >= p.cfg.Monsters
}

func (p *Processor) allMonsterFloorsCleared(player *domain.PlayerState) bool {
	for _, killed := range player.MonstersKilledByFloor {
		if killed < p.cfg.Monsters {
			return false
		}
	}

	return true
}

func (p *Processor) handleRegister(ev config.Event) {
	player := p.getOrCreatePlayer(ev.PlayerID)
	if player.Finished {
		p.handleImpossibleMove(ev)
		return
	}

	player.Registered = true
	p.output = append(p.output, ev)
}

func (p *Processor) handleEnterDungeon(ev config.Event) {
	player := p.getOrCreatePlayer(ev.PlayerID)

	if !player.Registered {
		p.handleDisqual(ev)
		return
	}
	if player.Finished || player.InDungeon || ev.Time.Before(p.cfg.OpenAt) || p.isDungeonClosed(ev) {
		p.handleImpossibleMove(ev)
		return
	}

	player.InDungeon = true
	player.Floor = 1
	player.EnteredAt = ev.Time
	if p.isMonsterFloor(player.Floor) {
		player.FloorEnterAt = ev.Time
	}

	p.output = append(p.output, ev)
}

func (p *Processor) handleKilledMonster(ev config.Event) {
	player := p.getOrCreatePlayer(ev.PlayerID)

	if !p.canHandleDungeonEvent(ev, player) {
		return
	}
	if !p.isMonsterFloor(player.Floor) {
		p.handleImpossibleMove(ev)
		return
	}

	floorIndex := player.Floor - 1
	if player.MonstersKilledByFloor[floorIndex] >= p.cfg.Monsters {
		p.handleImpossibleMove(ev)
		return
	}

	player.MonstersKilledByFloor[floorIndex]++
	if player.MonstersKilledByFloor[floorIndex] == p.cfg.Monsters {
		player.FloorClearDurations = append(player.FloorClearDurations, ev.Time.Sub(player.FloorEnterAt))
	}

	p.output = append(p.output, ev)
}

func (p *Processor) handleWentNextFloor(ev config.Event) {
	player := p.getOrCreatePlayer(ev.PlayerID)

	if !p.canHandleDungeonEvent(ev, player) {
		return
	}
	if player.Floor >= p.cfg.Floors || (p.isMonsterFloor(player.Floor) && !p.isCurrentFloorCleared(player)) {
		p.handleImpossibleMove(ev)
		return
	}

	player.Floor++
	if p.isMonsterFloor(player.Floor) && !p.isCurrentFloorCleared(player) {
		player.FloorEnterAt = ev.Time
	}

	p.output = append(p.output, ev)
}

func (p *Processor) handleWentPreviousFloor(ev config.Event) {
	player := p.getOrCreatePlayer(ev.PlayerID)

	if !p.canHandleDungeonEvent(ev, player) {
		return
	}
	if player.Floor <= 1 {
		p.handleImpossibleMove(ev)
		return
	}

	player.Floor--
	if p.isMonsterFloor(player.Floor) && !p.isCurrentFloorCleared(player) {
		player.FloorEnterAt = ev.Time
	}

	p.output = append(p.output, ev)
}

func (p *Processor) handleEnteredBossFloor(ev config.Event) {
	player := p.getOrCreatePlayer(ev.PlayerID)

	if !p.canHandleDungeonEvent(ev, player) {
		return
	}
	if player.Floor != p.cfg.Floors || !p.allMonsterFloorsCleared(player) || player.BossKilled {
		p.handleImpossibleMove(ev)
		return
	}

	player.BossEnterAt = ev.Time
	p.output = append(p.output, ev)
}

func (p *Processor) handleKilledBoss(ev config.Event) {
	player := p.getOrCreatePlayer(ev.PlayerID)

	if !p.canHandleDungeonEvent(ev, player) {
		return
	}
	if player.Floor != p.cfg.Floors || player.BossEnterAt.IsZero() || player.BossKilled {
		p.handleImpossibleMove(ev)
		return
	}

	player.BossKilled = true
	player.BossKillDuration = ev.Time.Sub(player.BossEnterAt)
	if p.allMonsterFloorsCleared(player) {
		player.Status = domain.StatusSuccess
	}

	p.output = append(p.output, ev)
}

func (p *Processor) handleLeftDungeon(ev config.Event) {
	player := p.getOrCreatePlayer(ev.PlayerID)

	if !p.canHandleDungeonEvent(ev, player) {
		return
	}

	player.InDungeon = false
	player.Finished = true
	player.FinishedAt = ev.Time
	if player.BossKilled && p.allMonsterFloorsCleared(player) {
		player.Status = domain.StatusSuccess
	} else {
		player.Status = domain.StatusFail
	}

	p.output = append(p.output, ev)
}

func (p *Processor) handleCannotContinue(ev config.Event) {
	player := p.getOrCreatePlayer(ev.PlayerID)

	if !player.Registered {
		p.handleDisqual(ev)
		return
	}
	if player.Finished {
		p.handleImpossibleMove(ev)
		return
	}

	player.InDungeon = false
	player.Finished = true
	player.FinishedAt = ev.Time
	player.Status = domain.StatusDisqual

	p.output = append(p.output, ev)
}

func (p *Processor) handleRestoredHealth(ev config.Event) {
	player := p.getOrCreatePlayer(ev.PlayerID)

	if !p.canHandleDungeonEvent(ev, player) {
		return
	}

	health, err := strconv.Atoi(ev.Extra)
	if err != nil || health < 0 {
		p.handleImpossibleMove(ev)
		return
	}

	player.HP += health
	if player.HP > 100 {
		player.HP = 100
	}

	p.output = append(p.output, ev)
}

func (p *Processor) handleReceivedDamage(ev config.Event) {
	player := p.getOrCreatePlayer(ev.PlayerID)

	if !p.canHandleDungeonEvent(ev, player) {
		return
	}

	damage, err := strconv.Atoi(ev.Extra)
	if err != nil || damage < 0 {
		p.handleImpossibleMove(ev)
		return
	}

	player.HP -= damage
	if player.HP < 0 {
		player.HP = 0
	}

	p.output = append(p.output, ev)
	if player.HP == 0 {
		player.InDungeon = false
		player.Finished = true
		player.FinishedAt = ev.Time
		player.Status = domain.StatusFail
		p.output = append(p.output, config.Event{
			Time:     ev.Time,
			PlayerID: ev.PlayerID,
			ID:       config.EventDead,
		})
	}
}

func (p *Processor) handleImpossibleMove(ev config.Event) {
	p.output = append(p.output, config.Event{
		Time:     ev.Time,
		PlayerID: ev.PlayerID,
		ID:       config.EventImpossibleMove,
		Extra:    fmt.Sprint(ev.ID),
	})
}

func (p *Processor) handleDisqual(ev config.Event) {
	player := p.getOrCreatePlayer(ev.PlayerID)
	if player.Status == domain.StatusDisqual {
		return
	}

	player.InDungeon = false
	player.Finished = true
	player.FinishedAt = ev.Time
	player.Status = domain.StatusDisqual
	p.output = append(p.output, config.Event{
		Time:     ev.Time,
		PlayerID: ev.PlayerID,
		ID:       config.EventDisqualified,
	})
}