package tool

//提供了一个非常简单的雪花ID生成器和解析器。
import (
	"github.com/bwmarrin/snowflake"
	"time"
)

var snowFlakeID *snowflake.Node

func InitSnowFlake(startTime string, machineID int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}
	snowflake.Epoch = st.UnixNano() / 1000000
	snowFlakeID, err = snowflake.NewNode(machineID)
	return
}

func GenID() int64 {
	return snowFlakeID.Generate().Int64()
}
