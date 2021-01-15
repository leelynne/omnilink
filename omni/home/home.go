package home

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/leelynne/omnilink/omni"
)

type Home struct {
	ModelNumber int
	ModelName   string
	Version     string
	PhoneNumber string
	features    []omni.SystemFeature

	mu           sync.Mutex
	thermostats  []Thermostat
	latestStatus Status
	troubles     []omni.SystemTrouble
}

func (h *Home) Features() []Feature {
	h.mu.Lock()
	defer h.mu.Unlock()

	features := []Feature{}
	for _, omniFeature := range h.features {
		features = append(features, Feature{int(omniFeature), omniFeature.String()})
	}

	return features
}
func (h *Home) Thermostats() ([]Thermostat, error) {
	return nil, nil
}

func (h *Home) String() string {
	buf := bytes.Buffer{}

	buf.WriteString(fmt.Sprintf("ModelNumber: %d\n", h.ModelNumber))
	buf.WriteString(fmt.Sprintf("ModelName: %s\n", h.ModelName))
	buf.WriteString(fmt.Sprintf("Version: %s\n", h.Version))
	buf.WriteString(fmt.Sprintf("Phone: %s\n", h.PhoneNumber))
	buf.WriteString(fmt.Sprintf("Features: %s\n", h.Features()))
	buf.WriteString(fmt.Sprintf("Status: %s\n", h.latestStatus))

	return buf.String()
}

type Feature struct {
	Type int
	Name string
}

func (f Feature) String() string {
	return f.Name
}

type Thermostat struct {
}

type Status struct {
	DateSet bool
	Date    time.Time
	Sunrise time.Time
	Sunset  time.Time
	Battery uint8
}

func (s Status) String() string {
	if s.DateSet {
		return fmt.Sprintf("Date: %s\nBattery: %d\n", s.Date, s.Battery)
	}
	return fmt.Sprintf("Date: Not Set\nBattery: %d\n", s.Battery)
}
