package model

import "time"

type Community struct {
	CommunityId   int       `db:"community_id" json:"communityId,string"`
	CommunityName string    `db:"community_name" json:"communityName"`
	Introduction  string    `db:"introduction" json:"introduction,omitempty"`
	CreateTime    time.Time `db:"create_time" json:"createTime,omitempty"`
}
