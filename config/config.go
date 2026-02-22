package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ExitCodes handles both single int and list of ints in YAML.
// e.g: exitcodes: 0  OR  exitcodes: [0, 2]
type ExitCodes []int

func (e *ExitCodes) UnmarshalYAML(value *yaml.Node) error {
	// Try as a single int first
	var single int
	if err := value.Decode(&single); err == nil {
		*e = ExitCodes{single}
		return nil
	}
	// Try as a list of ints
	var list []int
	if err := value.Decode(&list); err == nil {
		*e = ExitCodes(list)
		return nil
	}
	return fmt.Errorf("exitcodes: expected int or list of ints")
}

// Program represents a single supervised program from the config file.
type Program struct {
	Cmd          string            `yaml:"cmd"`
	NumProcs     int               `yaml:"numprocs"`
	AutoStart    bool              `yaml:"autostart"`
	AutoRestart  string            `yaml:"autorestart"` // always, never, unexpected
	ExitCodes    ExitCodes         `yaml:"exitcodes"`
	StartTime    int               `yaml:"starttime"`
	StartRetries int               `yaml:"startretries"`
	StopSignal   string            `yaml:"stopsignal"`
	StopTime     int               `yaml:"stoptime"`
	Stdout       string            `yaml:"stdout"`
	Stderr       string            `yaml:"stderr"`
	Env          map[string]string `yaml:"env"`
	WorkingDir   string            `yaml:"workingdir"`
	Umask        string            `yaml:"umask"`
}

// Config is the top-level structure of the YAML config file.
type Config struct {
	Programs map[string]Program `yaml:"programs"`
}

// Load reads and parses a YAML config file.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("cannot parse config file: %w", err)
	}

	return &cfg, nil
}
