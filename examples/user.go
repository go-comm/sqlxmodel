package examples

type Base struct {
	ID         string `json:"id,omitempty" db:"id"`
	CreateTime int64  `json:"createtime,omitempty" db:"createtime"`
	Creater    string `json:"creater,omitempty" db:"creater"`
	ModifyTime int64  `json:"modifytime,omitempty" db:"modifytime"`
	Modifier   string `json:"modifier,omitempty" db:"modifier"`
	Version    int    `json:"version,omitempty" db:"version"`
	Defunct    bool   `json:"defunct,omitempty" db:"defunct"`
	Deleted    bool   `json:"deleted,omitempty" db:"deleted"`
}

func (Base) PrimaryKey() string {
	return "id"
}

type User struct {
	Base
	Name  string `json:"name,omitempty" db:"name"`
	Email string `json:"email,omitempty" db:"email"`

	RoleID int64 `json:"role_id,omitempty" db:"role_id"`
	Role   *Role `json:"role,omitempty"`
}

func (User) TableName() string {
	return "t_user"
}
