package schema

type Config struct {
	Model Model
}

type Model map[string]ModelInfo
type ModelInfo map[string]Info

type Info struct {
	Type      string `yaml:"type"`
	FixedName bool   `yaml:"fixedName"`
}
