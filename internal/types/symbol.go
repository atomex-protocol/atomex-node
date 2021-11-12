package types

// Symbol -
type Symbol struct {
	Name     string `yaml:"name" validate:"required"`
	BaseKey  string `yaml:"base" validate:"required"`
	QuoteKey string `yaml:"quote" validate:"required"`

	Base  Asset `yaml:"-" validate:"-"`
	Quote Asset `yaml:"-" validate:"-"`
}
