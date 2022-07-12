package model

import "time"

type PostMsg struct {
	PostID      int64     `db:"post_id" json:"postID,string" `
	Title       string    `db:"title" json:"title" binding:"required"`
	Content     string    `db:"content" json:"content" binding:"required"`
	AuthorID    int64     `db:"author_id" json:"authorID,string"`
	CommunityID int       `db:"community_id" json:"communityID,string" binding:"required"`
	Status      int       `db:"status" json:"status,omitempty"`
	CreateTime  time.Time `db:"create_time" json:"createTime"`
	UpdateTime  time.Time `db:"update_time" json:"updateTime"`
}

type PostDetail struct {
	*PostMsg
	*Community
	*User
	VoteData int64 `json:"voteData"`
}
