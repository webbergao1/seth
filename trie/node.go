package trie

import "seth/common"

//NodeType define type NodeType
type NodeType int8

const (
	// NodeTypeNil node type nil, for error node
	NodeTypeNil NodeType = 0
	// NodeTypeExtend is extend node type
	NodeTypeExtend NodeType = 1
	// NodeTypeLeaf is leaf node type
	NodeTypeLeaf NodeType = 2
	// NodeTypeBranch is branch node type
	NodeTypeBranch NodeType = 3
)

// Noder interface for node
type Noder interface {
	Hash() common.Hash
	IsDirty() bool
}

// Node is trie node struct
type Node struct {
	hash  []byte // hash of the node
	dirty bool   // is the node dirty,need to write to database
}

// ExtendNode is extend node struct.for root,extend node
type ExtendNode struct {
	Node
	Key      []byte // for shared nibbles or key-end
	Nextnode Noder
}

// LeafNode is leaf node struct
type LeafNode struct {
	Node
	Key   []byte // for shared nibbles or key-end
	Value []byte // the value of leafnode
}

// BranchNode is node for branch
type BranchNode struct {
	Node
	Children [17]Noder
}

// hashNode is just used by nextnode of ExtendNode
// when it does not load real node from datbase
type hashNode []byte

// Hash return the hash of node
func (n hashNode) Hash() common.Hash { return common.BytesToHash(n) }

// IsDirty is node dirty
func (n hashNode) IsDirty() bool { return false }

// Hash return the hash of node
func (n Node) Hash() common.Hash {
	return common.BytesToHash(n.hash)
}

// IsDirty is node dirty
func (n Node) IsDirty() bool {
	return n.dirty
}
