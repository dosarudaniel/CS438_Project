package types

// Settings represents the main gossiper settings
type Settings struct {
	GossipAddr  string
	Name        string
	AntiEntropy int
	RouteTimer  int
}
