package controller

// LevelOrderIterator is an iterator for level-order (breadth-first) traversal.
type LevelOrderIterator struct {
	queue []*Node
}

// NewLevelOrderIterator creates a new iterator starting at the root.
func NewLevelOrderIterator(root *Node) *LevelOrderIterator {
	it := &LevelOrderIterator{}
	if root != nil {
		it.queue = append(it.queue, root)
	}
	return it
}

// HasNext returns true if there are more nodes to iterate.
func (it *LevelOrderIterator) HasNext() bool {
	return len(it.queue) > 0
}

// Next returns the next node's key in level order.
// It panics if there are no more nodes; in a production setting you might return an error instead.
func (it *LevelOrderIterator) Next() *Node {
	if !it.HasNext() {
		panic("No more elements")
	}
	// Dequeue the first node.
	node := it.queue[0]
	it.queue = it.queue[1:]
	// Enqueue the children.
	if node.Left != nil {
		it.queue = append(it.queue, node.Left)
	}
	if node.Right != nil {
		it.queue = append(it.queue, node.Right)
	}
	return node
}
