package settings

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf = new(AppConfig) //返回的值是一个指向该类型新分配的零值的指针

type AppConfig struct {
	Name      string `mapstructure:"name"`
	Mode      string `mapstructure:"mode"`
	Version   string `mapstructure:"version"`
	Port      int    `mapstructure:"port"`
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machineID"`

	*LogConfig   `mapstructure:"log"`
	*MysqlConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}
type MysqlConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Dbname       string `mapstructure:"dbname"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	Password string `mapstructure:"password"`
}

//viper使用配置管理
func Init(path string) (err error) {
	//方式1:直接指定配置文件路径,
	//相对路径,相对执行的可执行文件的相对路径
	viper.SetConfigFile(path)
	//相对路径:
	//viper.SetConfigFile("系统绝对路径")

	//方式2:指定配置文件和配置文件的位置,viper自行寻找可用配置文件
	//配置文件名不需要带后缀
	//配置文件位置可配置多个
	//viper.SetConfigName("config") //指定配置文件名称(不需要带后缀)
	//viper.AddConfigPath(".")      //指定查找配置文件的路径

	//基本上是配合远程配置中心使用,告诉viper当前的数据使用什么格式解析
	//viper.SetConfigType("yaml") //指定配置文件的类型(专用来从远程获取配置指定后缀)

	err = viper.ReadInConfig() //读取配置信息
	if err != nil {
		fmt.Printf("viper.ReadInConfig() failed ,err: %v \n", err)
		return
	}
	//把读取到的配置信息反序列化到Conf变量中
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed,err:", err)
	}
	//监控文件
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Printf("配置文件修改了")
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal err:", err)
		}
	})
	return
}
