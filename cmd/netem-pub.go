package main

import (
	"expvar"
	"fmt"
	"net/http"
	"time"

	"github.com/thomas-fossati/netem-pub/cmd/config"
	"github.com/thomas-fossati/netem-pub/netem"
)

type ifaceExpVars struct {
	PktCount     *expvar.Int
	PktDropped   *expvar.Int
	PktReordered *expvar.Int
	BytesCount   *expvar.Int
}

var netemExpVars map[string]ifaceExpVars

func updateExpVars(iface config.Interface, d *netem.NetemData) {
	v := netemExpVars[iface.Name]

	v.PktCount.Set(d.Total)
	v.PktDropped.Set(d.Dropped)
	v.PktReordered.Set(d.Reordered)
	v.BytesCount.Set(d.Bytes)
}

func initExpVars(cfg *config.Config) {
	netemExpVars = make(map[string]ifaceExpVars)

	for _, iface := range cfg.Interfaces {
		v := ifaceExpVars{
			PktCount:     expvar.NewInt(fmt.Sprintf("%s.packet.count", iface.Domain)),
			PktDropped:   expvar.NewInt(fmt.Sprintf("%s.packet.dropped", iface.Domain)),
			PktReordered: expvar.NewInt(fmt.Sprintf("%s.packet.reordered", iface.Domain)),
			BytesCount:   expvar.NewInt(fmt.Sprintf("%s.bytes.count", iface.Domain)),
		}

		netemExpVars[iface.Name] = v
	}
}

func netemPoll(cfg *config.Config) {
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

			updateExpVars(iface, netemData)

		}
		time.Sleep(cfg.PollIntervalMs * time.Millisecond)

	}
}

func NetemPubSvc(cfg *config.Config) {
	initExpVars(cfg)
	go netemPoll(cfg)
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.HTTPPort), nil)
}
