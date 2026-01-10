package entity

type User struct {
	ID        int64
	Username  string
	Email     string
	Role      string
	Avatar    string // 微信头像
	CreatedAt int64
	UpdatedAt int64
}
