package main

import (
	"context"
	"fmt"
	"log"
	"middleware/config"
	"sync"

	"github.com/IBM/sarama"
)

var (
	defaultProducer sarama.SyncProducer
	defaultConsumer sarama.Consumer
)

type Kafka struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

func InitKafka(ctx context.Context, wg *sync.WaitGroup) {
	var kafakCfg Kafka
	config.CfgKafka.Unmarshal(&kafakCfg)

	var err error
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	saramaConfig.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	saramaConfig.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回
	addr := fmt.Sprintf("%s:%s", kafakCfg.Host, kafakCfg.Port)

	// 连接kafka, 初始化生产者
	defaultProducer, err = sarama.NewSyncProducer([]string{addr}, saramaConfig)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to start producer, err:%v", err))
		return
	}

	// 连接kafka, 初始化消费者
	defaultConsumer, err = sarama.NewConsumer([]string{addr}, nil)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to start consumer, err:%v", err))
		return
	}

	log.Fatal("init kafka consumer success")
}

func GetKafkaConsumer() sarama.Consumer {
	return defaultConsumer
}

func GetKafkaProducer() sarama.SyncProducer {
	return defaultProducer
}
