package models

// LoginUser represents the user login credentials
type LoginUser struct {
	Username string `json:"username" binding:"required,min=3,max=32"` // 用户名必填，长度3-32
	Password string `json:"password" binding:"required,min=3"`        // 密码必填，最小长度6
}

// Validate 验证用户输入
func (u *LoginUser) Validate() bool {
	return len(u.Username) >= 3 && len(u.Username) <= 32 && len(u.Password) >= 3
}
