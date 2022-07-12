package redis

const (
	Prefix             = "web_app:"    // 项目key前缀
	KeyPostTimeZSet    = "post:time"   // zset;贴子及发帖时间
	KeyPostScoreZSet   = "post:score"  // zset;贴子及投票的分数
	KeyPostVotedZSetPF = "post:voted:" // zset;记录用户及投票类型;参数是post id
	KeyCommunitySetPF  = "community:"  // set;保存每个分区下帖子的id
)

// getRedisKey 拼接一下redis命名空间
func getRedisKey(key string) string {
	return Prefix + key
}
