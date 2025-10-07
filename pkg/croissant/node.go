// node.go
package croissant

// Node represents a node in the Croissant metadata structure.
type Node interface {
	GetName() string
	GetID() string
	GetParent() Node
	SetParent(Node)
	Validate(*Issues)
}

// BaseNode implements common functionality for all nodes.
type BaseNode struct {
	ID     string `json:"@id,omitempty"`
	Name   string `json:"name,omitempty"`
	parent Node
}

func (n *BaseNode) GetName() string {
	return n.Name
}

func (n *BaseNode) GetID() string {
	return n.ID
}

func (n *BaseNode) GetParent() Node {
	return n.parent
}

func (n *BaseNode) SetParent(parent Node) {
	n.parent = parent
}
