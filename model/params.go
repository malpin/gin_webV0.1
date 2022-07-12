package model

const (
	OrderTime  = "time"
	OrderScore = "score"
)

//请求参数的结构体再此定义

// ParamsSignUp 注册的请求参数
type ParamsSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// ParamsLogin 登录的请求参数
type ParamsLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ParamVoteData struct {
	//UserID 从请求中获取当前的用户
	PostID    int64 `json:"post_id,string" binding:"required"`
	Direction int   `json:"direction,string" binding:"oneof=1 0 -1"`
}

type ParamPostList struct {
	Page        int64  `json:"page" form:"page"`
	Size        int64  `json:"size" form:"size"`
	Order       string `json:"order" form:"order"`
	CommunityID int    `json:"communityID" form:"communityID"`
}
