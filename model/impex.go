package model

// ImpEx is the data structure for importing and exporting data to a and from
// the application
type ImpEx struct {
	Config  *Config  `json:"config"`
	Entries Entries  `json:"entries"`
	Weights []Weight `json:"weights"`
}
