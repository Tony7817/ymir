package id

import (
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// https://en.wikipedia.org/wiki/Snowflake_ID
//| 1 bit  |   41 bits   |  10 bits  |  12 bits  |
//| 固定 0  | 时间戳差值 | 机器 ID  | 序列号    |

const (
	epoch         = int64(1577836800000) // 自定义基准时间：2020-01-01 00:00:00 UTC
	machineIDBits = 10                   // 机器 ID 所占位数
	sequenceBits  = 12                   // 序列号所占位数

	maxMachineID = -1 ^ (-1 << machineIDBits) // 最大机器 ID
	maxSequence  = -1 ^ (-1 << sequenceBits)  // 最大序列号

	timeShift      = machineIDBits + sequenceBits
	machineIDShift = sequenceBits
)

var SF *Snowflake

func init() {
	SF = NewSnowFlake()
}

func NewSnowFlake() *Snowflake {
	var marchineIdRaw = os.Getenv("MACHINE_ID")
	if marchineIdRaw == "" {
		panic("MACHINE_ID is not set")
	}
	mId, err := strconv.ParseInt(marchineIdRaw, 10, 64)
	if err != nil {
		panic("MACHINE_ID is not a number")
	}
	return &Snowflake{
		machineID: mId,
		lastTime:  0,
		sequence:  0,
	}
}

type Snowflake struct {
	machineID int64
	lastTime  int64
	sequence  int64
	mutex     sync.Mutex
}

// GenerateID 生成唯一 ID
func (s *Snowflake) GenerateID() int64 {
	s.mutex.Lock()

	now := time.Now().UnixMilli()
	if now < s.lastTime {
		s.lastTime += 1
		logx.Errorf("clock is moving backwards. Rejecting requests until %ds", s.lastTime-now)
	}

	if now == s.lastTime {
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 {
			// 当前毫秒内序列号已用完，等待下一毫秒
			for now <= s.lastTime {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		// 时间进到下一毫秒，序列号重置
		s.sequence = 0
	}
	s.lastTime = now

	s.mutex.Unlock()

	return (now-epoch)<<timeShift | (s.machineID << machineIDShift) | s.sequence
}
