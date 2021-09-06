package common

import (
	"fmt"
	"time"

	"github.com/projectdiscovery/clistats"
)

type Progress interface {
	Init(Count int64)
	IncrementRequests()
}

type StatsTicker struct {
	active       bool
	tickDuration time.Duration
	stats        clistats.StatisticsClient
}

func NewStatsTicker(duration int, active bool) (Progress, error) {
	var tickDuration time.Duration
	if active {
		tickDuration = time.Duration(duration) * time.Second
	} else {
		tickDuration = -1
	}

	progress := &StatsTicker{}

	stats, err := clistats.New()
	if err != nil {
		return nil, err
	}
	progress.active = active
	progress.stats = stats
	progress.tickDuration = tickDuration

	return progress, nil
}

func (p *StatsTicker) IncrementRequests() {
	p.stats.IncrementCounter("requests", 1)
}

func (p *StatsTicker) Init(Count int64) {
	p.stats.AddStatic("hosts", Count)
	p.stats.AddStatic("startedAt", time.Now())
	p.stats.AddCounter("requests", uint64(0))

	if p.active {
		if err := p.stats.Start(printCallback, p.tickDuration); err != nil {
			fmt.Println(err)
		}
	}
}

func printCallback(stats clistats.StatisticsClient) {
	requests, _ := stats.GetCounter("requests")
	Bar.Set(int(requests))
}
