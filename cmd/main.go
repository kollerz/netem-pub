package main

import (
	"expvar"
	"net/http"
	"time"

	"github.com/thomas-fossati/netem-pub/netem"
)

var pktCount = expvar.NewInt("packet.count")
var pktDropped = expvar.NewInt("packet.dropped")
var pktReordered = expvar.NewInt("packet.reordered")
var bytesCount = expvar.NewInt("bytes.count")

func NetemPoll() {
	ifaces := []string{"eth0"} // TODO(tho) cfg
	for {
		for _, iface := range ifaces {
			out, err := netem.Fetch(iface)
			if err != nil {
				continue
			}

			netemData, err := netem.Parse(out)
			if err != nil {
				continue
			}

			updateExpVars(netemData)

		}
		time.Sleep(250 * time.Millisecond) // TODO(tho) cfg

	}
}

func updateExpVars(d *netem.NetemData) {
	pktCount.Set(d.Total)
	pktDropped.Set(d.Dropped)
	pktReordered.Set(d.Reordered)
	bytesCount.Set(d.Bytes)
}

func main() {
	go NetemPoll()
	http.ListenAndServe(":8080", nil)
}
