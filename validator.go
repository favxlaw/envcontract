package envcontract

type FieldContract struct {
	Name       string
	EnvKey     string
	Required   bool
	HasDefault bool
	Default    string
	Kind       string
}
