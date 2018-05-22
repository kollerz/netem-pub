package service

import (
	"expvar"
	"fmt"
	"sync"
	"time"

	"github.com/thomas-fossati/netem-pub/hping"
	"github.com/thomas-fossati/netem-pub/netem"
	"github.com/thomas-fossati/netem-pub/netemd/config"
)

type ifaceExpVars struct {
	PktCount     *expvar.Int
	PktDropped   *expvar.Int
	PktReordered *expvar.Int
	BytesCount   *expvar.Int
	ForwardDelay *expvar.Int
	ReverseDelay *expvar.Int
}

type expVars struct {
	Map map[string]ifaceExpVars
	Mtx sync.Mutex
}

var ev expVars

func updateNetemExpVars(iface config.Interface, d *netem.NetemData) {
	ev.Mtx.Lock()
	defer ev.Mtx.Unlock()

	v := ev.Map[iface.Name]

	v.PktCount.Set(d.Total)
	v.PktDropped.Set(d.Dropped)
	v.PktReordered.Set(d.Reordered)
	v.BytesCount.Set(d.Bytes)
}

func updateHpingExpVars(iface config.Interface, d *hping.HpingData) {
	ev.Mtx.Lock()
	defer ev.Mtx.Unlock()

	v := ev.Map[iface.Name]

	v.ForwardDelay.Set(d.ForwardDelay)
	v.ReverseDelay.Set(d.ReverseDelay)
}

func initExpVars(cfg *config.Config) {
	ev.Map = make(map[string]ifaceExpVars)

	for _, iface := range cfg.Interfaces {
		v := ifaceExpVars{
			PktCount:     expvar.NewInt(fmt.Sprintf("%s.packet.count", iface.Tag)),
			PktDropped:   expvar.NewInt(fmt.Sprintf("%s.packet.dropped", iface.Tag)),
			PktReordered: expvar.NewInt(fmt.Sprintf("%s.packet.reordered", iface.Tag)),
			BytesCount:   expvar.NewInt(fmt.Sprintf("%s.bytes.count", iface.Tag)),
			ForwardDelay: expvar.NewInt(fmt.Sprintf("%s.delay.forward", iface.Tag)),
			ReverseDelay: expvar.NewInt(fmt.Sprintf("%s.delay.reverse", iface.Tag)),
		}

		ev.Map[iface.Name] = v
	}
}

func netemPoller(cfg *config.Config) {
	for {
		for _, iface := range cfg.Interfaces {
			out, err := netem.Fetch(iface.Name)
			if err != nil {
				continue
			}

			netemData, err := netem.Parse(out)
			if err != nil {
				continue
			}

			updateNetemExpVars(iface, netemData)

		}
		time.Sleep(cfg.PollIntervalMs * time.Millisecond)
	}
}

func hpingPoller(cfg *config.Config) {
	for {
		for _, iface := range cfg.Interfaces {
			out, err := hping.Fetch(iface.PingHost)
			if err != nil {
				continue
			}

			hpingData, err := hping.Parse(out)
			if err != nil {
				continue
			}

			updateHpingExpVars(iface, hpingData)

		}
		time.Sleep(cfg.PollIntervalMs * time.Millisecond)
	}
}

func NetemPub(cfg *config.Config, noPing bool) {
	initExpVars(cfg)
	go netemPoller(cfg)
	if !noPing {
		go hpingPoller(cfg)
	}
}
