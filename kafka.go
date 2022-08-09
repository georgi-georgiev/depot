package tst

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

func GetOffsets(broker string, topic string, partition int) (int64, int64) {

	conn, err := kafka.DialLeader(context.Background(), "tcp", broker, topic, partition)
	if err != nil {
		panic(err)
	}

	currentOffset, _ := conn.Offset()
	lastOffset, err := conn.ReadLastOffset()
	if err != nil {
		panic(err)
	}

	log.Printf("current offset: %v", currentOffset)
	log.Printf("highest offset: %v", lastOffset)

	return currentOffset, lastOffset
}

func ConsumeAll() {
	// conf := GetConfig()

	// for _, topic := range conf.Kafka.Topics {
	// 	Consume(topic, -1)
	// }
}

func Consume(broker string, topic string, partition int) [][]byte {
	log.Printf("consume message topic: %s", topic)
	log.Printf("Parition:%v", partition)

	currentOffset, highestOffset := GetOffsets(broker, topic, partition)

	var messages [][]byte = make([][]byte, 0)

	kafkaConfig := kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		//GroupID:   conf.Kafka.Group,
		Partition: partition,
		//StartOffset: highestOffset - 1,
	}
	r := kafka.NewReader(kafkaConfig)
	//r.SetOffset(highestOffset - 1)

	for {
		if currentOffset >= highestOffset-1 {
			break
		}

		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			panic(err)
		}

		currentOffset = msg.Offset
		highestOffset = msg.HighWaterMark

		log.Printf("consume message offset: %v", msg.Offset)
		log.Printf("consume message high water mar: %v", msg.HighWaterMark)

		messages = append(messages, msg.Value)

		r.CommitMessages(context.Background(), msg)
	}

	r.Close()

	return messages
}

//TODO: remove placeholder if a topic has unique key
func ConsumeExact(broker string, topic string, key string, value string) []byte {
	time.Sleep(5 * time.Second)

	partition := FindPartition(topic, value)

	var messages [][]byte
	for i := 0; i < partition; i++ {
		messages = append(messages, Consume(broker, topic, i)...)
	}

	placeholder := "\"" + key + "\":" + value

	for i := len(messages) - 1; i >= 0; i-- {
		message := messages[i]
		if strings.Contains(string(message), placeholder) {
			return message
		}
	}

	panic("could not find message with placeholder:" + placeholder)
}

var topicPartition map[string]int = map[string]int{
	"Staging.Internal.QrState.Update.V1": 6,
}

func FindPartition(topic string, key string) int {
	partition, ok := topicPartition[topic]

	if !ok {
		return 1
	} else {
		return partition
	}

	// conf := GetConfig()
	// kafkaPartitions, err := kafka.LookupPartitions(context.Background(), "tcp", conf.Kafka.Broker, topic)
	// if err != nil {
	// 	panic(err)
	// }

	// var partitions []int
	// for _, kafkaPartition := range kafkaPartitions {
	// 	partitions = append(partitions, kafkaPartition.ID)
	// }

	// log.Println(partitions)

	// sort.Ints(partitions)

	// hash := murmur2([]byte(key))
	// log.Printf("message key hash:%v", hash)
	// log.Printf("topic partitions count:%v", len(kafkaPartitions))
	// idx := AbsInt(int(hash)) % len(kafkaPartitions)
	// log.Printf("topic partition index:%v", idx)
	// return partitions[idx]
}

func murmur2(data []byte) int32 {
	length := int32(len(data))
	seed := uint32(0x9747b28c)
	m := int32(0x5bd1e995)
	r := uint32(24)
	h := int32(seed ^ uint32(length))
	length4 := length / 4
	for i := int32(0); i < length4; i++ {
		i4 := i * 4
		k := int32(data[i4+0]&0xff) + (int32(data[i4+1]&0xff) << 8) + (int32(data[i4+2]&0xff) << 16) + (int32(data[i4+3]&0xff) << 24)
		k ^= int32(uint32(k) >> r)
		k *= m
		h *= m
		h ^= k
	}
	switch length % 4 {
	case 3:
		h ^= int32(data[(length & ^3)+2]&0xff) << 16
		fallthrough
	case 2:
		h ^= int32(data[(length & ^3)+1]&0xff) << 8
		fallthrough
	case 1:
		h ^= int32(data[length & ^3] & 0xff)
		h *= m
	}
	h ^= int32(uint32(h) >> 13)
	h *= m
	h ^= int32(uint32(h) >> 15)
	return h
}
