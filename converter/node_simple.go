package converter

type NodeSimple struct {
	*NodeBase
}

func NewNodeSimple(ctx *NodeContext, sheet Sheet, field *Field) *NodeSimple {
	n := &NodeSimple{
		NodeBase: NewNodeBase(ctx, sheet, field),
	}
	n.ParseNodes()
	return n
}

func (n *NodeSimple) Type() NodeType {
	return NodeTypeSimple
}

func (n *NodeSimple) ParseNodes() {

}
