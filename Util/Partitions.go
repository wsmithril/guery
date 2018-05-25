package Util

import (
	"fmt"
	"io"
)

type Partitions struct {
	Metadata *Metadata
	Rows     []*Row
	Index    int
}

func NewPartitions(md *Metadata) *Partitions {
	return &Partitions{
		Metadata: md,
		Rows:     []*Row{},
		Index:    0,
	}
}

func (self *Partitions) Read() (*Row, error) {
	if self.Index >= len(self.Rows) {
		return nil, io.EOF
	}
	self.Index++
	return self.Rows[self.Index-1], nil
}

func (self *Partitions) Write(row *Row) {
	self.Rows = append(self.Rows, row)
	self.Keys = row.Keys
}

func (self *Partitions) Reset() {
	self.Index = 0
}

func (self *Partitions) GetIndex(name string) int {
	if i, ok := self.Metadata.ColumnMap[name]; ok {
		return i
	}
	return -1
}

func (self *Partitions) GetKeyString() string {
	res := ""
	for _, key := range self.Keys {
		res += fmt.Sprintf("%v", key)
	}
	return res
}

func (self *Partitions) GetPartitionsNum() int {
	return len(self.Rows)
}