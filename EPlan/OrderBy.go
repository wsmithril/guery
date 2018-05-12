package EPlan

import (
	. "github.com/xitongsys/guery/Plan"
	"github.com/xitongsys/guery/Util"
	"github.com/xitongsys/guery/pb"
)

type EPlanOrderByNode struct {
	Location pb.Location
	Inputs   []pb.Location
	Output   pb.Location
	Metadata *Util.Metadata
}

func (self *EPlanOrderByNode) GetNodeType() EPlanNodeType {
	return EORDERBYNODE
}

func (self *EPlanOrderByNode) GetOutputs() []pb.Location {
	return []pb.Location{self.Output}
}

func (self *EPlanOrderByNode) GetLocation() pb.Location {
	return self.Location
}

func NewEPlanOrderByNode(node *PlanOrderByNode, inputs []pb.Location, output pb.Location) *EPlanOrderByNode {
	return &EPlanOrderByNode{
		Location: output,
		Inputs:   inputs,
		Output:   output,
		Metadata: node.GetMetadata(),
	}
}
