package trie

import (
	"bytes"
	"fmt"
	"io"
	"seth/common"
	"seth/database"
	"seth/rlp"
)

//Trie is a Merkle Patricia Trie
type Trie struct {
	db     database.Database
	root   Noder  // root node of the Trie
	prefix []byte //prefix of Trie node
}

// NewTrie new a trie tree
func NewTrie(root common.Hash, prefix []byte, db database.Database) (*Trie, error) {

	trie := &Trie{
		db:     db,
		prefix: prefix,
	}

	if (root != common.Hash{}) {
		rootnode, err := trie.loadNode(root[:])
		if err != nil {
			return nil, err
		}
		trie.root = rootnode
	}

	return trie, nil
}

// Update update [key,value] in the trie
func (t *Trie) Update(key, value []byte) {
	key = keybytesToHex(key)
	node, err := t.insert(t.root, key, value)
	if err == nil {
		t.root = node
	}
}

// Delete delete node with key in the trie
func (t *Trie) Delete(key []byte) {

}

// Get get the value by key
func (t *Trie) Get(key []byte) []byte {
	key = keybytesToHex(key)
	return t.get(t.root, key, 0)
}

// Hash return the hash of trie
func (t Trie) Hash() common.Hash {
	if t.root != nil {
		return t.root.Hash()
	}
	return common.Hash{}
}

// Commit commit the dirty node to database
func (t *Trie) Commit(batch database.Database) (common.Hash, error) {
	return common.Hash{}, nil
}

func (t *Trie) insert(node Noder, key []byte, value []byte) (Noder, error) {

	switch n := node.(type) {
	case *ExtendNode:
		fmt.Println(n.Key)
	case *LeafNode:
		matchlen := matchkeyLen(n.Key, key)
		if matchlen == len(n.Key) {
			n.Value = value
			n.dirty = true
			return n, nil
		}
		branchnode := &BranchNode{
			Node: Node{
				dirty: true,
			},
		}
		var err error
		branchnode.Children[n.Key[matchlen]], err = t.insert(nil, n.Key[matchlen+1:], n.Value)
		if err != nil {
			return nil, err
		}
		branchnode.Children[key[matchlen]], err = t.insert(nil, key[matchlen+1:], value)
		if err != nil {
			return nil, err
		}
		if matchlen == 0 {
			return branchnode, nil
		}

		return &ExtendNode{
			Node: Node{
				dirty: true,
			},
			Key:      key[:matchlen],
			Nextnode: branchnode,
		}, nil

	case nil:
		newnode := &LeafNode{
			Node: Node{
				dirty: true,
			},
			Key:   key,
			Value: value,
		}
		return newnode, nil
	}
	return nil, nil
}

// loadNode get node from memory cache or database
func (t *Trie) loadNode(hash []byte) (Noder, error) {
	//TODO need cache nodes
	key := append(t.prefix, hash...)
	val, err := t.db.Get(key)
	if err != nil || val == nil {
		return nil, err
	}
	return decodeNode(val)
}

// decodeNode decode node from buf byte
func decodeNode(value []byte) (Noder, error) {
	if len(value) == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	vals, _, err := rlp.SplitList(value)
	if err != nil {
		return nil, err
	}
	switch n, _ := rlp.CountValues(vals); n {
	case 2:
	case 17:
	default:
		return nil, nil
	}
	return nil, nil
}

func (t *Trie) get(node Noder, key []byte, pos int) (value []byte) {
	switch n := (node).(type) {
	case nil:
		return nil
	case *ExtendNode:
		if len(key)-pos < len(n.Key) || !bytes.Equal(n.Key, key[pos:pos+len(n.Key)]) {
			return nil
		}
		return t.get(n.Nextnode, key, pos+len(n.Key))
	case hashNode:
		child, err := t.loadNode(n)
		if err != nil {
			return nil
		}
		return t.get(child, key, pos)
	case *LeafNode:
		if len(key)-pos < len(n.Key) || !bytes.Equal(n.Key, key[pos:pos+len(n.Key)]) {
			// key not found in trie
			return nil
		}
		return n.Value
	case *BranchNode:
		return t.get(n.Children[key[pos]], key, pos+1)
	default:
		panic(fmt.Sprintf("invalid node: %v", node))
	}
}

func keybytesToHex(str []byte) []byte {
	l := len(str)*2 + 1
	var nibbles = make([]byte, l)
	for i, b := range str {
		nibbles[i*2] = b / 16
		nibbles[i*2+1] = b % 16
	}
	nibbles[l-1] = 16
	return nibbles
}

func matchkeyLen(a, b []byte) int {
	var i, length = 0, len(a)
	if len(b) < length {
		length = len(b)
	}
	for ; i < length; i++ {
		if a[i] != b[i] {
			break
		}
	}
	return i
}
