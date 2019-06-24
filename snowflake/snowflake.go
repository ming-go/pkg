/*
	Reference:
		[1] https://discordapp.com/developers/docs/reference#snowflakes
		[2] https://github.com/twitter-archive/snowflake
*/

package snowflake

import (
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ming-go/pkg/mtime"
)

type snowflake struct {
	sync.Mutex
	datacenterId  int64
	workerId      int64
	sequence      int64
	lastTimestamp int64
}

const epoch = 1561344367000
const workerIdBits = 5
const datacenterIdBits = 5
const sequenceBits = 12

const workerIdShift = sequenceBits
const datacenterIdShift = sequenceBits + workerIdBits
const timestampLeftShift = sequenceBits + workerIdBits + datacenterIdBits

var sequenceMask int64 = -1 ^ (-1 << sequenceBits)
var maxWorkerId int64 = -1 ^ (-1 << workerIdBits)
var maxDatacenterId int64 = -1 ^ (-1 << datacenterIdBits)

const BIT_BETWEEN_21_TO_17 = 0x3E0000 // DataCenter ID
const BIT_BETWEEN_16_TO_12 = 0x1F000  // Process ID
const BIT_BETWEEN_11_TO_0 = 0xFFF     // Increment

func New(workerId int64, datacenterId int64) (*snowflake, error) {
	if workerId > maxWorkerId || workerId < 0 {
		return nil, errors.New(fmt.Sprintf("worker Id can't be greater than %d or less than 0", maxWorkerId))
	}

	if datacenterId > maxDatacenterId || datacenterId < 0 {
		return nil, errors.New(fmt.Sprintf("datacenter Id can't be greater than %d or less than 0", maxDatacenterId))
	}

	return &snowflake{
		datacenterId: datacenterId,
		workerId:     workerId,
	}, nil
}

func (sf *snowflake) NextId() (int64, error) {
	timestamp := sf.timeGen()
	sf.Lock()
	defer sf.Unlock()

	if timestamp < sf.lastTimestamp {
		return -1, errors.New(fmt.Sprintf("Clock moved backwards. Refusing to generate id for %d milliseconds", (sf.lastTimestamp - timestamp)))
	}

	if timestamp == sf.lastTimestamp {
		sf.sequence = (sf.sequence + 1) & sequenceMask
		if sf.sequence == 0 {
			timestamp = sf.nextMillis()
		}
	} else {
		sf.sequence = 0
	}

	return ((sf.timeGen() - epoch) << timestampLeftShift) |
		(sf.datacenterId << datacenterIdShift) |
		(sf.workerId << workerIdShift) |
		sf.sequence, nil
}

func (sf *snowflake) NextBase62Id() (string, error) {
	id, err := sf.NextId()
	if err != nil {
		return "", err
	}

	return big.NewInt(id).Text(62), nil
}

func (sf *snowflake) RetrievalTimestamp() {

}

func (sf *snowflake) nextMillis() int64 {
	timestamp := sf.timeGen()
	for timestamp < sf.lastTimestamp {
		timestamp = sf.timeGen()
	}

	return timestamp
}

func (sf *snowflake) timeGen() int64 {
	return mtime.UnixMilli(time.Now())
}
