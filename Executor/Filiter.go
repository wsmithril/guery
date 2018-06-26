package Executor

import (
	"fmt"
	"io"

	"github.com/vmihailenco/msgpack"
	"github.com/xitongsys/guery/Config"
	"github.com/xitongsys/guery/EPlan"
	"github.com/xitongsys/guery/Logger"
	"github.com/xitongsys/guery/Metadata"
	"github.com/xitongsys/guery/Split"
	"github.com/xitongsys/guery/Util"
	"github.com/xitongsys/guery/pb"
)

func (self *Executor) SetInstructionFilter(instruction *pb.Instruction) (err error) {
	var enode EPlan.EPlanFilterNode
	if err = msgpack.Unmarshal(instruction.EncodedEPlanNodeBytes, &enode); err != nil {
		return err
	}
	self.Instruction = instruction
	self.EPlanNode = &enode
	self.InputLocations = []*pb.Location{&enode.Input}
	self.OutputLocations = []*pb.Location{&enode.Output}
	return nil
}

func (self *Executor) RunFilter() (err error) {
	defer self.Clear()

	if self.Instruction == nil {
		return fmt.Errorf("No Instruction")
	}
	enode := self.EPlanNode.(*EPlan.EPlanFilterNode)

	md := &Metadata.Metadata{}
	reader := self.Readers[0]
	writer := self.Writers[0]
	if err = Util.ReadObject(reader, md); err != nil {
		return err
	}

	//write metadata
	if err = Util.WriteObject(writer, md); err != nil {
		return err
	}

	rbReader := Split.NewSplitBuffer(md, reader, nil)
	rbWriter := Split.NewSplitBuffer(md, nil, writer)

	//write
	jobs := make(chan *Split.Split)
	done := make(chan bool)

	for i := 0; i < int(Config.Conf.Runtime.ParallelNumber); i++ {
		go func() {
			defer func() {
				done <- true
			}()

			for {
				sp, ok := <-jobs
				if ok {
					for i := 0; i < sp.GetRowsNumber(); i++ {
						flag := true
						for _, booleanExpression := range enode.BooleanExpressions {
							if ok, err := booleanExpression.Result(sp, i); !ok.(bool) && err == nil {
								flag = false
								break
							} else if err != nil {
								flag = false
								break
							}
						}

						if flag {
							if err = rbWriter.Write(sp, i); err != nil {
								continue //should add err handler
							}
						}
					}

				} else {
					break
				}
			}
		}()
	}

	var sp *Split.Split
	for err == nil {
		sp, err = rbReader.ReadSplit()
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			break
		}
		jobs <- sp
	}
	close(jobs)
	for i := 0; i < int(Config.Conf.Runtime.ParallelNumber); i++ {
		<-done
	}

	if err = rbWriter.Flush(); err != nil {
		return err
	}

	Logger.Infof("RunFilter finished")
	return err
}
