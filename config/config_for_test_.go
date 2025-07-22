package conf

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

//这段代码的目的是为单元测试提供一个灵活的配置加载方式。它通过定义一个接口 (IReader)，使得测试时可以不依赖真实的配置文件，而是传入一个模拟的（mock）配置内容，这让测试更加独立和可控。
// 这个文件是为了方便写的test文件来读取config

type IReader interface {
	readConfig() ([]byte, error)
}

type ConfigReader struct {
	FileName string
}

// 'reader' implementing the Interface
// Function to read from actual file
func (r *ConfigReader) readConfig() ([]byte, error) {
	file, err := ioutil.ReadFile(r.FileName)

	if err != nil {
		log.Fatal(err)
	}
	return file, err
}

func InitConfigForTest(reader IReader) {
	file, err := reader.readConfig()
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(file, &Config)
	if err != nil {
		panic(err)
	}
}
