package engine

type Output struct {
	Entity  string            `json:"entity"`
	ID      string            `json:"id"`
	Date    string            `json:"date"`
	Type    string            `json:"type"`
	Value   string            `json:"value"`
	Details map[string]string `json:"details,omitempty"`
}
