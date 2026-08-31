package label

type Label struct {
	Address  string `gorm:"primaryKey" json:"address"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Source   string `json:"source"`
}
