package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

)

type Config struct {
	Floors 		int 		`json:"Floors"`
	Monsters 	int			`json:"Monsters"`
	OpenAtStr	string		`json:"OpenAt"`
	OpenAt		time.Time
	DurationInt	int			`json:"Duration"`
	Duration	time.Duration
}

func LoadConfig(Path string) (*Config, error){
	file, err := os.Open(Path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open file %s: %v", Path, err)
		return nil, err;
	}
	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		fmt.Fprintf(os.Stderr, "failed to decode %s: %v", Path, err)
		return nil, err
	}

	if err := cfg.ValidateConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to validate config: %v", err)
		return nil, err
	}

	if cfg.Floors == 0 || cfg.Monsters == 0{
		err :=  fmt.Errorf("invalid floors count")
		fmt.Fprintf(os.Stderr, "failed to create dungeon: %v", err)
		return nil, err
	}

	return &cfg, nil
}

func (cfg *Config) ValidateConfig() error {
	s := strings.Trim(cfg.OpenAtStr, "\"")

	if strings.Count(s, ":") != 2 {
		return fmt.Errorf("invalid time format: expected HH:MM:SS, got %s", s)
	}

	t, err := time.Parse("15:04:05", s)
	if err != nil {
		return fmt.Errorf("failed to parse time: %w", err)
	}

	cfg.OpenAt = t
	cfg.Duration = time.Hour * time.Duration(cfg.DurationInt)

	return nil
}