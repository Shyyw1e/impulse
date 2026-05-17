package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type EventID int

//ingoing events
const (
	EventRegistered        EventID = 1
	EventEnteredDungeon    EventID = 2
	EventKilledMonster     EventID = 3
	EventWentNextFloor     EventID = 4
	EventWentPreviousFloor EventID = 5
	EventEnteredBossFloor  EventID = 6
	EventKilledBoss        EventID = 7
	EventLeftDungeon       EventID = 8
	EventCannotContinue    EventID = 9
	EventRestoredHealth    EventID = 10
	EventReceivedDamage    EventID = 11
)

//outgoing events
const (
	EventDisqualified   EventID = 31
	EventDead           EventID = 32
	EventImpossibleMove EventID = 33
)

type Event struct {
	Time     time.Time
	PlayerID int
	ID       EventID
	Extra    string
}

func LoadEvents(path string) ([]Event, error) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open %s: %v", path, err)
		return nil, err
	}
	defer file.Close()

	var events []Event

	scanner := bufio.NewScanner(file)
	linenumber := 0
	for scanner.Scan() {
		linenumber++
		line := strings.Fields(scanner.Text())
		if len(line) == 0 {
			continue
		}
		if len(line) < 3 {
			return nil, fmt.Errorf("line %d: invalid event format", linenumber)
		}

		event, err := parseEventLine(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse line %d: %v", linenumber, err)
			return nil, fmt.Errorf("line %d: %w", linenumber, err)
		}

		events = append(events, event)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func parseEventLine(line []string) (Event, error) {
	if len(line) < 3 {
		return Event{}, fmt.Errorf("invalid event format")
	}

	pID, err := strconv.Atoi(line[1])
	if err != nil {
		return Event{}, err
	}
	id, err := strconv.Atoi(line[2])
	if err != nil {
		return Event{}, err
	}
	event := Event{
		PlayerID: pID,
		ID:       EventID(id),
	}
	rawTime := line[0]
	rawTime = strings.Trim(rawTime, "[")
	rawTime = strings.Trim(rawTime, "]")
	if strings.Count(rawTime, ":") != 2 {
		return Event{}, fmt.Errorf("invalid event time")
	}
	t, err := time.Parse("15:04:05", rawTime)
	if err != nil {
		return event, fmt.Errorf("failed to parse time: %w", err)
	}
	event.Time = t
	if len(line) > 3 {
		event.Extra = strings.Join(line[3:], " ")
	}

	return event, nil
}
