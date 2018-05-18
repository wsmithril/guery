package Optimizer

import (
	"sort"

	"github.com/xitongsys/guery/Plan"
)

func FilterColumns(node Plan.PlanNode, columns []string) error {
	if node == nil {
		return nil
	}
	switch node.(type) {
	case *Plan.PlanLimitNode:
		fallthrough
	case *Plan.PlanUnionNode:
		fallthrough
	case *Plan.PlanCombineNode:
		indexes := []int{}
		md := node.GetMetadata()
		for _, c := range columns {
			if index, err := md.GetIndexByName(c); err == nil {
				indexes = append(indexes, index)
			}
		}
		sort.Ints(indexes)

		inputs := node.GetInputs()
		mdis := []*Util.Metadata{}
		for _, input := range inputs {
			mdis = append(mdis, input.GetMetadata())
		}

		intputNum := len(node.GetInputs)
		columnsForInput := make([][]string, len(inputNum))

		i, indexNum := 0, mdis[0].GetColumnNumber()
		for _, index := range indexes {
			if index < indexNum {
				indexForInput := index - (indexNum - mdis[i].GetColumnNumber())
				cname := mdis[i].Columns[indexForInput].Name
				columnsForInput[i] = append(columnsForInput, cname)
			} else {
				i++
				indexNum += mdis[i].GetColumnNumber()
			}
		}
		for i, input := range inputs {
			err := FilterColumns(input, columnsForInput[i])
			if err != nil {
				return err
			}
		}
	case *Plan.PlanFiliterNode:
		nodea := node.(*Plan.PlanFiliterNode)
		columnsForInput, err := nodea.BooleanExpression.GetColumns()
		if err != nil {
			return err
		}
		return FilterColumns(nodea.Input, columnsForInput)

	case *Plan.PlanGroupByNode:
		nodea := node.(*Plan.PlanGroupByNode)
		columnsForInput := []string{}
		for _, ele := range nodea.GroupingElement {
			cs, err := ele.GetColumns()
			if err != nil {
				return err
			}
			columnsForInput = append(columnsForInput, cs...)
		}
		cs, err := nodea.Having.GetColumns()
		if err != nil {
			return err
		}
		columnsForInput = append(columnsForInput, cs...)
		return FilterColumns(nodea.Input, columnsForInput)

	case *Plan.PlanJoinNode:
		nodea := node.(*Plan.PlanJoinNode)
		columns := []string{}

	case *Plan.PlanOrderByNode:
		nodea := node.(*Plan.PlanOrderByNode)
		columnsForInput := []string{}
		for _, item := range nodea.SortItems {
			cs, err := item.GetColumns()
			if err != nil {
				return err
			}
			columnsForInput = append(columnsForInput, cs...)
		}
		return FilterColumns(nodea.Input, columnsForInput)

	case *Plan.PlanSelectNode:
		nodea := node.(*Plan.PlanJoinNode)
		columnsForInput := []string{}
		for _, item := range self.SelectItems {
			cs, err := item.GetColumns()
			if err != nil {
				return err
			}
			columnsForInput = append(columnsForInput, cs...)
		}
		return FilterColumns(nodea.Input, columnsForInput)

	case *Plan.PlanScanNode:
		nodea := node.(*Plan.PlanScanNode)

	case *Plan.PlanRenameNode: //already use deleteRenameNode
		return nil
	default:
		return fmt.Errorf("Unknown PlanNode type")
	}

	return nil
}
