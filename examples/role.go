package examples

type Role struct {
	ID   int64  `json:"id,omitempty" db:"id"`
	Name string `json:"name,omitempty" db:"name"`
}

func (Role) PrimaryKey() string {
	return "id"
}

func (Role) TableName() string {
	return "t_role"
}
