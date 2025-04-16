package apptypes

type RoleID int

const (
	Guest RoleID = iota // 游客
	User                // 用户
	Admin               // 管理员
)
