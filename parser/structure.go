package parser

type Application struct {
	Download struct {
		URL          string `json:"URL"`
		IsCompressed bool   `json:"IsCompressed"`
	} `json:"Download"`
	Execution struct {
		FileLocation string   `json:"FileLocation"`
		Arguments    []string `json:"Arguments"`
		Requirements []string `json:"Requirements"`
		RunAsAdmin   bool     `json:"RunAsAdmin"`
	} `json:"Execution"`
	Metadata struct {
		Name        string   `json:"Name"`
		Description string   `json:"Description"`
		Author      string   `json:"Author"`
		Tags        []string `json:"Tags"`
	} `json:"Metadata"`
	Compatibility struct {
		OS            []string `json:"OS"`
		Architectures []string `json:"Architectures"`
	} `json:"Compatibility"`
}
