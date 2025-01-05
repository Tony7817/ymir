package id

import (
	"sync"
	"testing"
	"time"
)

var SF *Snowflake

func init() {
	SF = NewSnowFlake()
}

func TestSnowFlake(t *testing.T) {
	id := SF.GenerateID()
	t.Log(id)
}

func TestGenerateID(t *testing.T) {
	// 测试生成唯一 ID
	id1 := SF.GenerateID()
	id2 := SF.GenerateID()
	if id1 == id2 {
		t.Error("expected unique IDs, got duplicate IDs")
	}

	// 测试生成 ID 的时间戳部分
	timestamp1 := (id1 >> timeShift) + epoch
	timestamp2 := (id2 >> timeShift) + epoch
	if timestamp1 > timestamp2 {
		t.Errorf("expected timestamp1 <= timestamp2, got %d > %d", timestamp1, timestamp2)
	}
}

func TestClockBackwards(t *testing.T) {
	// 模拟正常时间
	now := time.Now().UnixMilli()
	SF.lastTime = now + 10 // 模拟 lastTime 比当前时间大

	// 测试时钟回拨情况
	SF.GenerateID() // 应触发 log
}

func TestSequenceOverflow(t *testing.T) {
	for i := 0; i < 1000; i++ {

		// 模拟同一毫秒内生成多个 ID
		SF.lastTime = time.Now().UnixMilli()
		SF.sequence = maxSequence // 模拟序列号已达最大值

		now := SF.lastTime
		// 此次生成应该进入下一毫秒
		id := SF.GenerateID()
		timestamp := (id >> timeShift) + epoch
		if timestamp <= now {
			t.Errorf("expected timestamp > lastTime, got %d <= %d", timestamp, SF.lastTime)
		}
	}
}

func TestConcurrentGenerateID(t *testing.T) {
	var wg sync.WaitGroup
	idSet := sync.Map{}

	// 并发生成 ID
	workerCount := 100
	idsPerWorker := 1000
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < idsPerWorker; j++ {
				id := SF.GenerateID()
				// 检查 ID 是否重复
				if _, exists := idSet.LoadOrStore(id, struct{}{}); exists {
					t.Errorf("duplicate ID detected: %d", id)
				}
			}
		}()
	}
	wg.Wait()
}

func TestEpoch(t *testing.T) {

	// 检查生成的 ID 时间戳是否从 epoch 开始计算
	id := SF.GenerateID()
	timestamp := (id >> timeShift) + epoch
	now := time.Now().UnixMilli()

	if timestamp > now {
		t.Errorf("expected timestamp <= now, got %d > %d", timestamp, now)
	}
}
