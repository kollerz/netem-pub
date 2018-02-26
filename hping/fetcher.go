package hping

import "os/exec"

func Fetch(addr string) (string, error) {
	out, err := exec.Command("/usr/sbin/hping", "-c", "1", "--icmp-ts", addr).Output()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
