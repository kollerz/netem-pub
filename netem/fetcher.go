package netem

import "os/exec"

func Fetch(iface string) (string, error) {
	out, err := exec.Command("/sbin/tc", "-s", "qdisc", "show", "dev", iface).Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
