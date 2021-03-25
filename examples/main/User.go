package main

type User struct {
	UID        string `json:"uid,omitempty" db:"uid"`
	Name       string `json:"name,omitempty" db:"name"`
	Email      string `json:"email,omitempty" db:"email"`
	CreateTime int64  `json:"createtime,omitempty" db:"createtime"`
	Creater    string `json:"creater,omitempty" db:"creater"`
	ModifyTime int64  `json:"modifytime,omitempty" db:"modifytime"`
	Modifier   string `json:"modifier,omitempty" db:"modifier"`
	Version    int    `json:"version,omitempty" db:"version"`
	Defunct    bool   `json:"defunct,omitempty" db:"defunct"`
	Deleted    bool   `json:"deleted,omitempty" db:"deleted"`
}

func (u User) TableName() string {
	return "t_user"
}

func (u User) PrimaryKey() string {
	return "uid"
}
