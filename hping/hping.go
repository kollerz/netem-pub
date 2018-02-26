package hping

type HpingFetcher interface {
	Fetch(iface string) (string, error)
}

type HpingParser interface {
	Parse(text string) (*HpingData, error)
}
