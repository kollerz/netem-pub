package config

import (
	"fmt"
	"net"
	"time"
)

type Config struct {
	PollIntervalMs time.Duration `mapstructure:"poll_interval_ms"`
	HTTPPort       int           `mapstructure:"http_port"`
	Interfaces     []Interface   `mapstructure:"interfaces"`
}

type Interface struct {
	Tag    string `mapstructure:"tag"`
	IPAddr string `mapstructure:"ipaddr"`
	Name   string // gets resolved from IPAddr on the fly
}

func NewConfig() *Config {
	return &Config{
		PollIntervalMs: 1000,
		HTTPPort:       8888,
	}
}

func (c *Config) Remap() error {
	for _, iface := range c.Interfaces {
		tmp, err := getInterfaceNameFromIP(iface.IPAddr)
		if err != nil {
			return err
		}
		iface.Name = tmp
	}
	return nil
}

func getInterfaceNameFromIP(ipaddr string) (string, error) {
	needle := net.ParseIP(ipaddr)
	if needle == nil {
		return "", fmt.Errorf("invalid IP address %s", ipaddr)
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 ||
			iface.Flags&net.FlagLoopback != 0 {
			continue // interface down or loopback
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range addrs {
			var ip net.IP

			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip != nil && ip.Equal(needle) {
				return iface.Name, nil
			}
		}
	}

	return "", fmt.Errorf("address %s could not be mapped to an existing interface", ipaddr)
}
