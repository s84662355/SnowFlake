package SnowFlake

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	WorkerIdBitsMax           uint8 = 14            //服务器id bit 最大长度
	WorkerIdBitsMin           uint8 = 8             //服务器id bit 最小长度
	NumberBitsMax             uint8 = 14            //自增id bit 最大长度
	NumberBitsMin             uint8 = 8             //自增id  bit 最小长度
	WorkerIdBitsAndNumberBits uint8 = 22            //服务器id bit +  自增id  bit
	EpochMin                  int64 = 1666195200000 //最小起始时间
	TimeStampBits             uint8 = 41            //时间戳 bit 长度
	CharacterBits             uint8 = 1
)

type ID struct {
	TimeStamp int64
	WorkId    int64
	Number    int64
	Id        int64
}

type Options struct {
	WorkerIdBits uint8 ///服务器id bit 长度
	NumberBits   uint8 ///自增id bit 长度
	Epoch        int64 //起始时间戳 毫秒
	WorkId       int64 //机器id
}

type App struct {
	workerIdBits  uint8
	numberBits    uint8
	numberMax     int64
	workerIdMax   int64
	timeShift     uint8
	workerIdShift uint8
	workId        int64
	epoch         int64
	mu            sync.Mutex
	timeStamp     int64
	number        int64
}

func New(o *Options) (*App, error) {
	now := time.Now().UnixNano() / 1e6

	if o.Epoch > now {
		return nil, errors.New("Epoch 参数 不能大于当前时间戳")
	}

	if o.Epoch < EpochMin {
		return nil, errors.New(fmt.Sprint("Epoch 参数 不能小于", EpochMin))
	}

	if o.WorkerIdBits > WorkerIdBitsMax {
		return nil, errors.New(fmt.Sprint("WorkerIdBits 参数 不能大于", WorkerIdBitsMax))
	}

	if o.WorkerIdBits < WorkerIdBitsMin {
		return nil, errors.New(fmt.Sprint("WorkerIdBits 参数 不能小于", WorkerIdBitsMin))
	}

	if o.NumberBits > NumberBitsMax {
		return nil, errors.New(fmt.Sprint("NumberBits 参数 不能大于", NumberBitsMax))
	}

	if o.NumberBits < NumberBitsMin {
		return nil, errors.New(fmt.Sprint("NumberBits 参数 不能小于", NumberBitsMin))
	}

	if o.WorkerIdBits+o.NumberBits != WorkerIdBitsAndNumberBits {
		return nil, errors.New(fmt.Sprint("WorkerIdBits + NumberBits 必须等于", WorkerIdBitsAndNumberBits))
	}

	r := &App{}

	r.workerIdBits = o.WorkerIdBits
	r.numberBits = o.NumberBits
	r.epoch = o.Epoch
	r.timeShift = WorkerIdBitsAndNumberBits
	r.workerIdShift = o.NumberBits

	r.numberMax = -1 ^ (-1 << r.numberBits)
	r.workerIdMax = -1 ^ (-1 << r.workerIdBits)

	if o.WorkId > r.workerIdMax || o.WorkId < 0 {
		return nil, errors.New(fmt.Sprint("WorkId 参数 必须大于0小于", r.workerIdMax))
	}

	r.workId = o.WorkId

	return r, nil
}

//生成id
func (r *App) NextId() *ID {
	r.mu.Lock()
	defer r.mu.Unlock()
	//当前时间的毫秒时间戳
	now := time.Now().UnixNano() / 1e6
	//如果时间戳与当前时间相同，则增加序列号
	if r.timeStamp == now {
		r.number++
		//如果序列号超过了最大值，则更新时间戳
		if r.number > r.numberMax {
			for now <= r.timeStamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else { //如果时间戳与当前时间不同，则直接更新时间戳
		r.number = 0
		r.timeStamp = now
	}
	//ID由时间戳、机器编码、序列号组成
	id := (now-r.epoch)<<r.timeShift | (r.workId << r.workerIdShift) | (r.number)

	idData := &ID{
		TimeStamp: now,
		WorkId:    r.workId,
		Number:    r.number,
		Id:        id,
	}
	return idData
}

//解析id
func (r *App) DecodeID(id int64) *ID {
	var (
		TimeStamp int64
		WorkId    int64
		Number    int64
	)
	TimeStamp = r.epoch + (id >> WorkerIdBitsAndNumberBits)
	WorkId = id << (TimeStampBits + CharacterBits) >> (TimeStampBits + CharacterBits + r.workerIdShift)
	Number = id << (TimeStampBits + CharacterBits + r.workerIdBits) >> (TimeStampBits + CharacterBits + r.workerIdBits)
	return &ID{
		TimeStamp: TimeStamp,
		WorkId:    WorkId,
		Number:    Number,
		Id:        id,
	}
}
