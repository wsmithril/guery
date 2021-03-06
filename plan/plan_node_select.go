package plan

import (
	"fmt"
	"strings"

	"github.com/xitongsys/guery/config"
	"github.com/xitongsys/guery/gtype"
	"github.com/xitongsys/guery/metadata"
	"github.com/xitongsys/guery/parser"
)

type PlanSelectNode struct {
	Input         PlanNode
	Output        PlanNode
	Metadata      *metadata.Metadata
	SetQuantifier *gtype.QuantifierType
	SelectItems   []*SelectItemNode
	Having        *BooleanExpressionNode
	IsAggregate   bool
}

func NewPlanSelectNode(runtime *config.ConfigRuntime, input PlanNode, sq parser.ISetQuantifierContext, items []parser.ISelectItemContext, having parser.IBooleanExpressionContext) *PlanSelectNode {
	res := &PlanSelectNode{
		Input:         input,
		Metadata:      metadata.NewMetadata(),
		SetQuantifier: nil,
		SelectItems:   []*SelectItemNode{},
		Having:        nil,
	}
	if sq != nil {
		q := gtype.StrToQuantifierType(sq.GetText())
		res.SetQuantifier = &q
	}
	for i := 0; i < len(items); i++ {
		itemNode := NewSelectItemNode(runtime, items[i])
		res.SelectItems = append(res.SelectItems, itemNode)
		if itemNode.IsAggregate() {
			res.IsAggregate = true
		}
	}

	if having != nil {
		res.Having = NewBooleanExpressionNode(runtime, having)
		if res.Having.IsAggregate() {
			res.IsAggregate = true
		}
	}
	return res
}

func (self *PlanSelectNode) GetNodeType() PlanNodeType {
	return SELECTNODE
}

func (self *PlanSelectNode) GetInputs() []PlanNode {
	return []PlanNode{self.Input}
}

func (self *PlanSelectNode) SetInputs(inputs []PlanNode) {
	self.Input = inputs[0]
}

func (self *PlanSelectNode) GetOutput() PlanNode {
	return self.Output
}

func (self *PlanSelectNode) SetOutput(output PlanNode) {
	self.Output = output
}

func (self *PlanSelectNode) GetMetadata() *metadata.Metadata {
	return self.Metadata
}

func (self *PlanSelectNode) SetMetadata() error {
	if err := self.Input.SetMetadata(); err != nil {
		return err
	}
	md := self.Input.GetMetadata()
	colNames, colTypes := []string{}, []gtype.Type{}
	for _, item := range self.SelectItems {
		names, types, err := item.GetNamesAndTypes(md)
		if err != nil {
			return err
		}
		colNames = append(colNames, names...)
		colTypes = append(colTypes, types...)
	}

	if len(colNames) != len(colTypes) {
		return fmt.Errorf("length error")
	}
	self.Metadata = metadata.NewMetadata()
	for i, name := range colNames {
		t := colTypes[i]
		column := metadata.NewColumnMetadata(t, strings.Split(name, ".")...)
		self.Metadata.AppendColumn(column)
	}
	self.Metadata.Reset()

	return nil
}

func (self *PlanSelectNode) String() string {
	res := "PlanSelectNode {\n"
	res += "Input: " + self.Input.String() + "\n"
	res += "Metadata: " + fmt.Sprint(self.Metadata) + "\n"
	res += "SelectItems: " + fmt.Sprint(self.SelectItems) + "\n"
	res += "}\n"
	return res
}
