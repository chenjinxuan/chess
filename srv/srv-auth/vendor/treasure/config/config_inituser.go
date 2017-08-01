package config

type UserInit struct {
	Enable  bool            `json:"enable"`
	Balance int             `json:"balance"`
	Filter  *UserInitFilter `json:"filter"`
}

type UserInitFilter struct {
	Enable bool                `json:"enable"`
	White  map[string][]string `json:"white"`
}
