package engine

// MappingConfig representa a configuração para mapear a entrada para saída
type MappingConfig struct {
	Mappings []Mapping `json:"mappings"`
}

// Mapping define um mapeamento entre a entrada e a saída
type Mapping struct {
	Entity      string            `json:"entity"`
	InputKey    string            `json:"inputKey"`
	OutputType  string            `json:"outputType"`
	DetailKeys  []string          `json:"detailKeys"`
	TypeMapping map[string]string `json:"typeMapping"`
}
