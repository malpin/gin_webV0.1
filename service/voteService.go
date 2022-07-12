package service

import (
	"gin_web/dao/redis"
	"gin_web/model"
	"strconv"
)

// PostVote 文章点赞
func PostVote(userID int64, vota *model.ParamVoteData) error {
	//转换格式
	formatInt := strconv.FormatInt(userID, 10)
	s := strconv.FormatInt(vota.PostID, 10)
	f := float64(vota.Direction)

	err := redis.VoteForPost(formatInt, s, f)
	if err != nil {
		return err
	}
	return nil
}
