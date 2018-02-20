package netem

type NetemFetcher interface {
	Fetch(iface string) (string, error)
}

type NetemParser interface {
	Parse(text string) (*NetemData, error)
}
