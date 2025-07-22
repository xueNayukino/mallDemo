package conf

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

// 定义与 YAML 配置文件结构对应的 Go `struct`，并使用 `viper` 库来读取 `config.yaml` 文件，然后将内容解析到这些 `struct` 中。
var Config *Conf

// Conf 是所有配置的根结构体，对应 YAML 文件的顶层
type Conf struct {
	System        *System                 `yaml:"system"`
	Oss           *Oss                    `yaml:"oss"`
	MySql         map[string]*MySql       `yaml:"mysql"`
	Email         *Email                  `yaml:"email"`
	Redis         *Redis                  `yaml:"redis"`
	EncryptSecret *EncryptSecret          `yaml:"encryptSecret"`
	Cache         *Cache                  `yaml:"cache"`
	KafKa         map[string]*KafkaConfig `yaml:"kafKa"`
	RabbitMq      *RabbitMq               `yaml:"rabbitMq"`
	Es            *Es                     `yaml:"es"`
	PhotoPath     *LocalPhotoPath         `yaml:"photoPath"`
}

type RabbitMq struct {
	RabbitMQ         string `yaml:"rabbitMq"`
	RabbitMQUser     string `yaml:"rabbitMqUser"`
	RabbitMQPassWord string `yaml:"rabbitMqPassWord"`
	RabbitMQHost     string `yaml:"rabbitMqHost"`
	RabbitMQPort     string `yaml:"rabbitMqPort"`
}

type KafkaConfig struct {
	DisableConsumer bool   `yaml:"disableConsumer"`
	Debug           bool   `yaml:"debug"`
	Address         string `yaml:"address"`
	RequiredAck     int    `yaml:"requiredAck"`
	ReadTimeout     int64  `yaml:"readTimeout"`
	WriteTimeout    int64  `yaml:"writeTimeout"`
	MaxOpenRequests int    `yaml:"maxOpenRequests"`
	Partition       int    `yaml:"partition"`
}

// // 注意：这里的 AppEnv 字段在 YAML 中是 env，可能导致读取不到
type System struct {
	AppEnv      string `yaml:"appEnv"`
	Domain      string `yaml:"domain"`
	Version     string `yaml:"version"`
	HttpPort    string `yaml:"httpPort"`
	Host        string `yaml:"host"`
	UploadModel string `yaml:"uploadModel"`
}

type Oss struct {
	BucketName      string `yaml:"bucketName"`
	AccessKeyId     string `yaml:"accessKeyId"`
	AccessKeySecret string `yaml:"accessKeySecret"`
	Endpoint        string `yaml:"endPoint"`
	EndpointOut     string `yaml:"endpointOut"`
	QiNiuServer     string `yaml:"qiNiuServer"`
}

type MySql struct {
	Dialect  string `yaml:"dialect"`
	DbHost   string `yaml:"dbHost"`
	DbPort   string `yaml:"dbPort"`
	DbName   string `yaml:"dbName"`
	UserName string `yaml:"userName"`
	Password string `yaml:"password"`
	Charset  string `yaml:"charset"`
}

type Email struct {
	ValidEmail string `yaml:"validEmail"`
	SmtpHost   string `yaml:"smtpHost"`
	SmtpEmail  string `yaml:"smtpEmail"`
	SmtpPass   string `yaml:"smtpPass"`
}

type Redis struct {
	RedisHost     string `yaml:"redisHost"`
	RedisPort     string `yaml:"redisPort"`
	RedisUsername string `yaml:"redisUsername"`
	RedisPassword string `yaml:"redisPwd"`
	RedisDbName   int    `yaml:"redisDbName"`
	RedisNetwork  string `yaml:"redisNetwork"`
}

// EncryptSecret 加密的东西
type EncryptSecret struct {
	JwtSecret   string `yaml:"jwtSecret"`
	EmailSecret string `yaml:"emailSecret"`
	PhoneSecret string `yaml:"phoneSecret"`
	MoneySecret string `yaml:"moneySecret"`
}

type LocalPhotoPath struct {
	PhotoHost   string `yaml:"photoHost"`
	ProductPath string `yaml:"productPath"`
	AvatarPath  string `yaml:"avatarPath"`
}

type Cache struct {
	CacheType    string `yaml:"cacheType"`
	CacheExpires int64  `yaml:"cacheExpires"`
	CacheWarmUp  bool   `yaml:"cacheWarmUp"`
	CacheServer  string `yaml:"cacheServer"`
}

type Es struct {
	EsHost  string `yaml:"esHost"`
	EsPort  string `yaml:"esPort"`
	EsIndex string `yaml:"esIndex"`
}

func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("config")                    //文件名
	viper.SetConfigType("yaml")                      //文件类型
	viper.AddConfigPath(workDir + "/config/locales") // 会在此路径下寻找 config.yaml
	viper.AddConfigPath(workDir)                     // 也会在当前工作目录下寻找 config.yaml
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	// 将读取到的配置信息反序列化到全局的 Config 变量中
	err = viper.Unmarshal(&Config)
	if err != nil {
		panic(err)
	}
}

func GetExpiresTime() int64 { // GetExpiresTime 获取缓存过期时间的辅助函数
	if Config.Cache.CacheExpires == 0 {
		return int64(30 * time.Minute) // 默认 30min
	} //默认时间

	if Config.Cache.CacheExpires == -1 {
		return -1 // Redis.KeepTTL = -1
	} //-1永不过期
	// // 根据配置的值（单位：分钟）计算出最终的过期时间（单位：纳秒）
	return int64(time.Duration(Config.Cache.CacheExpires) * time.Minute)
}
