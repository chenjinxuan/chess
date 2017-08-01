package config

type TokenInfoCheck struct {
	Android TokenInfoCheckDetail `json:"android"`
	Ios     TokenInfoCheckDetail `json:"ios"`
}

type TokenInfoCheckDetail struct {
	Check bool `json:"check"`
	Min   int  `json:"min"`
	Max   int  `json:"max"`
}

type Base struct {
	Mod     string `json:"mod"`
	Key     string `json:"key"`
	ProdUrl string `json:"prod_url"`
}
