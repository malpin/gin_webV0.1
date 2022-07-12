package redis

import (
	"context"
	"fmt"
	"gin_web/Bean"
	"gin_web/model"
	"github.com/go-redis/redis/v9"
	"math"
	"strconv"
	"time"
)

var ctx, cancelFunc = context.WithTimeout(context.Background(), 5000*time.Second)

const voteTime = 60 * 60 * 24 * 180
const scorePerVote = 432

func VoteForPost(userID, postID string, direction float64) error {
	//根据帖子的id获取了发布时间
	// 1. 判断投票限制
	// 去redis取帖子发布时间
	creatTime := redisdb.ZScore(ctx, getRedisKey(KeyPostTimeZSet), postID).Val()

	//当前时间减去发帖时间,计算是否过期
	//判断一下是否在180天之内
	unix := time.Now().Unix()
	if float64(unix)-creatTime > voteTime {
		return Bean.VOTE_TIME_EXPIRED.MarkError
	}

	// 2. 更新贴子的分数
	//查询当前用户给当前帖子投票的信息
	val := redisdb.ZScore(ctx, getRedisKey(KeyPostVotedZSetPF+postID), userID).Val()
	// 更新：如果这一次投票的值和之前保存的值一致，就提示不允许重复投票
	if direction == val {
		fmt.Println("不允许重复投票运行")
		return Bean.ERR_VOTE_REPEATED.MarkError
	}
	var dir float64
	if direction > val {
		dir = 1
	} else {
		dir = -1
	}
	diff := math.Abs(val - direction) //计算两次投票的差值
	pipeline := redisdb.TxPipeline()  //开启事务
	//修改帖子的分数
	pipeline.ZIncrBy(ctx, getRedisKey(KeyPostScoreZSet), dir*diff*scorePerVote, postID)

	// 3. 记录用户为该贴子投票的数据
	if direction == 0 {
		//删除用户为帖子的投票数据
		pipeline.ZRem(ctx, getRedisKey(KeyPostVotedZSetPF+postID), userID)
	} else {
		//新增用户为帖子的投票数据
		pipeline.ZAdd(ctx, getRedisKey(KeyPostVotedZSetPF+postID), redis.Z{
			Score:  direction, // 赞成票还是反对票
			Member: userID,    //哪个用户
		})
	}
	_, err := pipeline.Exec(ctx) //提交事务
	return err
}

// 添加帖子的id 发布时间
func AddPortCreatTime(postID string, communityID int) error {
	pipeline := redisdb.TxPipeline() //事务
	// 帖子时间
	pipeline.ZAdd(ctx, getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	// 帖子分数
	pipeline.ZAdd(ctx, getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	//添加社区 id :帖子id
	ckey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(communityID))
	pipeline.SAdd(ctx, ckey, postID)
	_, err := pipeline.Exec(ctx)
	return err
}

//查找帖子 排序
func GetPostListOrder(p *model.ParamPostList) ([]string, error) {
	key := getRedisKey(KeyPostTimeZSet) //默认按照发帖时间查询
	//从redis获取id
	//根据用户请求中的order信息确定查询规则
	if p.Order == model.OrderScore { //按照分数排序
		key = getRedisKey(KeyPostScoreZSet)
	}
	//确定查询的索引起始点
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	// ZRevRange 按分数从大到小查询  ZRange 从小到大
	result, err := redisdb.ZRevRange(ctx, key, start, end).Result()
	return result, err
}

//查询帖子的点赞数量
func GetPostVoteData(ids []string) (data []int64, err error) {
	//发送多条命令
	pipeline := redisdb.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPF + id)
		pipeline.ZCount(ctx, key, "1", "1")
	}
	exec, err := pipeline.Exec(ctx)
	if err != nil {
		return
	}
	//获取返回的数值
	data = make([]int64, 0, len(ids))
	for _, i := range exec {
		val := i.(*redis.IntCmd).Val()
		data = append(data, val)
	}
	return
}

func GetPostListOrderByCommunity(p *model.ParamPostList) ([]string, error) {
	//使用zinterstore 把分区的帖子set 与帖子分数的zset 生成新的zset 联合查询
	//针对新的zset 按照之前的逻辑取数据
	tkey := getRedisKey(KeyPostTimeZSet) //默认按照发帖时间查询
	//从redis获取id
	//根据用户请求中的order信息确定查询规则
	if p.Order == model.OrderScore { //按照分数排序
		tkey = getRedisKey(KeyPostScoreZSet)
	}
	//利用缓存key减少使用zinterstore执行次数
	key := tkey + strconv.Itoa(p.CommunityID)
	ckey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(p.CommunityID))

	//判断key存在
	if redisdb.Exists(ctx, p.Order).Val() < 1 {
		pipeline := redisdb.TxPipeline()
		pipeline.ZInterStore(ctx, key, &redis.ZStore{
			Keys:      []string{ckey, tkey},
			Aggregate: "MAX",
		})
		pipeline.Expire(ctx, key, 600*time.Hour) //设置到期时间
		_, err := pipeline.Exec(ctx)
		if err != nil {
			return nil, err
		}
	}

	//确定查询的索引起始点
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	// ZRevRange 按分数从大到小查询  ZRange 从小到大
	result, err := redisdb.ZRevRange(ctx, key, start, end).Result()
	if err != nil {
		return nil, err
	}
	return result, err
}
