package code

import (
	"IDGeneratorService/config"
	"sync"
)

const DefaultEpoch int64 = 1564588800

//2018/1/1 for generate 19 digit
const MysqlEpochRemain int64 = 29326441

//2012/01/01 for generate 19 digit
const ManualEpochRemain int64 = 239710972

var once sync.Once
var generator *Generator

type Generator struct {
	node *Node
}

func InitGenerator(workID int64, IsMySQLAssigner bool) *Generator {
	once.Do(func() {
		var node *Node
		var err error
		if IsMySQLAssigner {
			node, err = NewNode(workID, MySQLAssignerOptions...)
			if err != nil {
				return
			}
		} else {
			node, err = NewNode(workID, ManualAssignOptions...)
			if err != nil {
				return
			}
		}
		generator = new(Generator)
		generator.node = node
	})
	return generator
}

func GetGenerator() *Generator {
	return generator
}

func (g *Generator) Next() int64 {
	return g.node.Generate().Int64()
}
func (g *Generator) Many(n int64) []int64 {
	r := make([]int64, 0)
	for i := 0; i < (int)(n); i++ {
		r = append(r, g.node.Generate().Int64())
	}
	return r
}

type (
	// Option allows specifying various settings
	Option func(*options)

	// options specify optional settings
	options struct {
		Epoch           int64
		NodeBits        uint8
		StepBits        uint8
		IsMySQLAssigner bool
	}
)

func EpochOpt(epoch int64) Option {
	return func(opts *options) {
		opts.Epoch = epoch
	}
}

func NodeBitsOpt(nodeBits uint8) Option {
	return func(opts *options) {
		opts.NodeBits = nodeBits
	}
}

func StepBitsOpt(stepBits uint8) Option {
	return func(opts *options) {
		opts.StepBits = stepBits
	}
}

func IsMysqlAssignerOpt(flag bool) Option {
	return func(opts *options) {
		opts.IsMySQLAssigner = flag
	}
}

var MySQLAssignerOptions = []Option{
	EpochOpt(ReSetEpoch(config.GetConfig().GetEpochUnixSecond(), true)),
	NodeBitsOpt(22),
	StepBitsOpt(13),
	IsMysqlAssignerOpt(true),
}

var ManualAssignOptions = []Option{
	EpochOpt(ReSetEpoch(config.GetConfig().GetEpochUnixSecond(), false)),
	NodeBitsOpt(10),
	StepBitsOpt(12),
	IsMysqlAssignerOpt(false),
}

func ReSetEpoch(customizeEpoch int64, isMySQLAssigner bool) int64 {
	if customizeEpoch == 0 {
		customizeEpoch = DefaultEpoch
	}
	return RemainEpochFor19Digits(customizeEpoch, isMySQLAssigner)
}

func RemainEpochFor19Digits(epoch int64, isMySQLAssigner bool) int64 {
	if isMySQLAssigner {
		return epoch - MysqlEpochRemain
	}
	return epoch - ManualEpochRemain
}
