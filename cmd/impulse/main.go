package main

import (
    "fmt"
    "os"

    "github.com/Shyyw1e/impulse/internal/config"
    "github.com/Shyyw1e/impulse/internal/output"
    "github.com/Shyyw1e/impulse/internal/usecase"
)

func main() {
    if len(os.Args) != 3 {
        fmt.Fprintf(os.Stderr, "usage: %s <config.json> <events>\n", os.Args[0])
        os.Exit(1)
    }

    configPath := os.Args[1]
    eventsPath := os.Args[2]

    cfg, err := config.LoadConfig(configPath)
    if err != nil {
        os.Exit(1)
    }

    events, err := config.LoadEvents(eventsPath)
    if err != nil {
        os.Exit(1)
    }

    processor := usecase.NewProcessor(cfg)
    processedEvents := processor.Process(events)

    for _, ev := range processedEvents {
        fmt.Println(output.FormatEvent(ev))
    }
}
