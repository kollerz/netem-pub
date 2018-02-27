package hping

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type HpingData struct {
	ForwardDelay int64
	ReverseDelay int64
}

func Parse(text string) (*HpingData, error) {
	var originateTs, receiveTs, transmitTs, tsrtt int64

	scanner := bufio.NewScanner(strings.NewReader(text))

	for lineNo := 0; scanner.Scan(); lineNo++ {
		var match bool
		line := scanner.Text()

		switch lineNo {
		case 0:
			match, _ = regexp.MatchString("HPING.*$", line)
			if !match {
				return nil, errors.New("Should start with HPING")
			}
		case 1:
			match, _ = regexp.MatchString("len=.*$", line)
			if !match {
				return nil, errors.New("Should start with len=")
			}
		case 2:
			match, _ = regexp.MatchString("^ICMP timestamp:.*$", line)
			if !match {
				return nil, errors.New("malformed hping sample (1)")
			}

			fmt.Sscanf(line, "ICMP timestamp: Originate=%d Receive=%d Transmit=%d",
				&originateTs, &receiveTs, &transmitTs)

		case 3:
			match, _ = regexp.MatchString("ICMP timestamp RTT.*$", line)
			if !match {
				return nil, errors.New("malformed hping sample (2)")
			}

			// extract tsrtt
			fmt.Sscanf(line, "ICMP timestamp RTT tsrtt=%d", &tsrtt)

			// do the maths to populate HpingData
			return &HpingData{
				ForwardDelay: receiveTs - originateTs,
				ReverseDelay: tsrtt + originateTs - transmitTs,
			}, nil
		}
	}

	return nil, errors.New("unexpected hping sample")
}
