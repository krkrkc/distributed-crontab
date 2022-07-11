package worker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	EtcdEndpoints            []string `json:"etcdEndpoints"`
	EtcdDialTimeout          int      `json:"etcdDialTimeout"`
	MongodbUri               string   `json:"mongodbUri"`
	MongodbConnectionTimeout int      `json:"mongodbConnectionTimeout"`
	JobLogBatchSize          int      `json:"jobLogBatchSize"`
	JobLogCommitTimeout      int      `json:"jobLogCommitTimeout"`
}

var (
	G_config *Config
)

func InitConfig(filename string) (err error) {
	var (
		content []byte
		conf    Config
	)

	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}

	if err = json.Unmarshal(content, &conf); err != nil {
		return
	}
	G_config = &conf
	fmt.Println(G_config)
	return
}
