package output

import (
	"fmt"
	"strings"
	"time"

	"github.com/Shyyw1e/impulse/internal/config"
)

func FormatEvents(events []config.Event) string {
	lines := make([]string, 0, len(events))
	for _, ev := range events {
		lines = append(lines, FormatEvent(ev))
	}

	return strings.Join(lines, "\n")
}

func FormatEvent(ev config.Event) string {
	prefix := fmt.Sprintf("[%s] Player [%d]", FormatTime(ev.Time), ev.PlayerID)

	switch ev.ID {
	case config.EventRegistered:
		return prefix + " registered"
	case config.EventEnteredDungeon:
		return prefix + " entered the dungeon"
	case config.EventKilledMonster:
		return prefix + " killed the monster"
	case config.EventWentNextFloor:
		return prefix + " went to the next floor"
	case config.EventWentPreviousFloor:
		return prefix + " went to the previous floor"
	case config.EventEnteredBossFloor:
		return prefix + " entered the boss's floor"
	case config.EventKilledBoss:
		return prefix + " killed the boss"
	case config.EventLeftDungeon:
		return prefix + " left the dungeon"
	case config.EventCannotContinue:
		return fmt.Sprintf("%s cannot continue due to %s", prefix, ev.Extra)
	case config.EventRestoredHealth:
		return fmt.Sprintf("%s has restored [%s] of health", prefix, ev.Extra)
	case config.EventReceivedDamage:
		return fmt.Sprintf("%s recieved [%s] of damage", prefix, ev.Extra)
	case config.EventDisqualified:
		return prefix + " is disqualified"
	case config.EventDead:
		return prefix + " is dead"
	case config.EventImpossibleMove:
		return fmt.Sprintf("%s makes imposible move [%s]", prefix, ev.Extra)
	default:
		return fmt.Sprintf("%s unknown event [%d]", prefix, ev.ID)
	}
}

func FormatTime(t time.Time) string {
	return t.Format("15:04:05")
}

func FormatDuration(d time.Duration) string {
	if d < 0 {
		d = 0
	}

	totalSeconds := int(d.Seconds())
	hours := totalSeconds / 3600
	minutes := totalSeconds % 3600 / 60
	seconds := totalSeconds % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
