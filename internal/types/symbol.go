package types

// Symbol -
type Symbol struct {
	Name     string `yaml:"name"`
	BaseKey  string `yaml:"base"`
	QuoteKey string `yaml:"quote"`

	Base  Asset `yaml:"-"`
	Quote Asset `yaml:"-"`
}
