package controller

import "fmt"

// Node represents a node in an AVL tree.
type Node struct {
	Key    int
	Left   *Node
	Right  *Node
	Height int
}

// NewNode creates a new Node with a given key.
func NewNode(key int) *Node {
	return &Node{
		Key:    key,
		Height: 1, // A new node is initially at height 1 (not 0).
	}
}

// height returns the height of a node. If node is nil, return 0.
func height(n *Node) int {
	if n == nil {
		return 0
	}
	return n.Height
}

// max returns the maximum of two integers.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// getBalance returns the balance factor of a node: height(left) - height(right).
// A balance factor in [-1, 0, +1] means it's balanced enough for AVL.
func getBalance(n *Node) int {
	if n == nil {
		return 0
	}
	return height(n.Left) - height(n.Right)
}

// rightRotate performs a right rotation on the subtree rooted at y and
// returns the new root (x).
//
//	  y           x
//	 / \         / \
//	x   T3  =>  T1  y
//
// / \             / \
// T1 T2           T2 T3
func rightRotate(y *Node) *Node {
	x := y.Left
	T2 := x.Right

	// Perform rotation
	x.Right = y
	y.Left = T2

	// Update heights
	y.Height = max(height(y.Left), height(y.Right)) + 1
	x.Height = max(height(x.Left), height(x.Right)) + 1

	// x is now root of the rotated subtree
	return x
}

// leftRotate performs a left rotation on the subtree rooted at x and
// returns the new root (y).
//
//	 x            y
//	/ \          / \
//
// T1  y   =>   x  T3
//
//	 / \      / \
//	T2 T3    T1 T2
func leftRotate(x *Node) *Node {
	y := x.Right
	T2 := y.Left

	// Perform rotation
	y.Left = x
	x.Right = T2

	// Update heights
	x.Height = max(height(x.Left), height(x.Right)) + 1
	y.Height = max(height(y.Left), height(y.Right)) + 1

	// y is now root of the rotated subtree
	return y
}

// Insert inserts a key into the AVL tree and returns the new root.
func Insert(node *Node, key int) *Node {
	// 1. Perform normal BST insertion
	if node == nil {
		return NewNode(key)
	}
	if key < node.Key {
		node.Left = Insert(node.Left, key)
	} else if key > node.Key {
		node.Right = Insert(node.Right, key)
	} else {
		// If you want to handle duplicates, decide whether to:
		// - disallow them,
		// - store a count, or
		// - insert to left/right consistently.
		// For simplicity, we'll ignore duplicates here (do nothing).
		return node
	}

	// 2. Update this node's height
	node.Height = 1 + max(height(node.Left), height(node.Right))

	// 3. Get the balance factor
	balance := getBalance(node)

	// 4. If the node is unbalanced, then try out the 4 cases

	// Case 1: Left Left
	if balance > 1 && key < node.Left.Key {
		return rightRotate(node)
	}

	// Case 2: Right Right
	if balance < -1 && key > node.Right.Key {
		return leftRotate(node)
	}

	// Case 3: Left Right
	if balance > 1 && key > node.Left.Key {
		node.Left = leftRotate(node.Left)
		return rightRotate(node)
	}

	// Case 4: Right Left
	if balance < -1 && key < node.Right.Key {
		node.Right = rightRotate(node.Right)
		return leftRotate(node)
	}

	// Return the (unchanged) node pointer
	return node
}

// minValueNode returns the node with the smallest key in the subtree.
func minValueNode(node *Node) *Node {
	current := node
	for current.Left != nil {
		current = current.Left
	}
	return current
}

// Remove deletes a key from the AVL tree and returns the new root.
func Remove(root *Node, key int) *Node {
	if root == nil {
		return root
	}

	// 1. Perform standard BST removal
	if key < root.Key {
		root.Left = Remove(root.Left, key)
	} else if key > root.Key {
		root.Right = Remove(root.Right, key)
	} else {
		// This is the node to be removed
		if root.Left == nil || root.Right == nil {
			// One child or no child
			var temp *Node
			if root.Left != nil {
				temp = root.Left
			} else {
				temp = root.Right
			}

			if temp == nil {
				// No child case
				root = nil
			} else {
				// One child case
				*root = *temp
			}
		} else {
			// Node with two children:
			// Get the inorder successor (smallest in the right subtree)
			temp := minValueNode(root.Right)
			root.Key = temp.Key
			// Remove the inorder successor
			root.Right = Remove(root.Right, temp.Key)
		}
	}

	// If the tree had only one node, then return
	if root == nil {
		return root
	}

	// 2. Update height of the current node
	root.Height = 1 + max(height(root.Left), height(root.Right))

	// 3. Check the balance factor
	balance := getBalance(root)

	// 4. Balance the node if needed

	// Left Left
	if balance > 1 && getBalance(root.Left) >= 0 {
		return rightRotate(root)
	}

	// Left Right
	if balance > 1 && getBalance(root.Left) < 0 {
		root.Left = leftRotate(root.Left)
		return rightRotate(root)
	}

	// Right Right
	if balance < -1 && getBalance(root.Right) <= 0 {
		return leftRotate(root)
	}

	// Right Left
	if balance < -1 && getBalance(root.Right) > 0 {
		root.Right = rightRotate(root.Right)
		return leftRotate(root)
	}

	return root
}

// Search returns true if key exists in the AVL tree, otherwise false.
func Search(root *Node, key int) bool {
	if root == nil {
		return false
	}
	if key < root.Key {
		return Search(root.Left, key)
	} else if key > root.Key {
		return Search(root.Right, key)
	}
	// key == root.Key
	return true
}

// InOrder traverses the AVL tree in ascending order and prints keys.
func InOrder(root *Node) {
	if root == nil {
		return
	}
	InOrder(root.Left)
	fmt.Printf("%d ", root.Key)
	InOrder(root.Right)
}

func LevelOrder(root *Node) {
	if root == nil {
		return
	}
	// Initialize the queue with the root node.
	queue := []*Node{root}
	for len(queue) > 0 {
		// Dequeue the first node.
		node := queue[0]
		queue = queue[1:]
		// Process the current node.
		fmt.Printf("%d ", node.Key)
		// Enqueue left child if it exists.
		if node.Left != nil {
			queue = append(queue, node.Left)
		}
		// Enqueue right child if it exists.
		if node.Right != nil {
			queue = append(queue, node.Right)
		}
	}
}

func demo() {
	var root *Node

	// Insert some values
	nums := []int{10, 20, 5, 6, 8, 15, 30, 25, 28}
	for _, n := range nums {
		root = Insert(root, n)
	}

	// Print the AVL tree in ascending order
	fmt.Print("InOrder traversal: ")
	InOrder(root)
	fmt.Println()

	// Check some searches
	fmt.Println("Search for 15:", Search(root, 15)) // true
	fmt.Println("Search for 7:", Search(root, 7))   // false

	// Remove a node
	fmt.Println("Removing 20...")
	root = Remove(root, 20)

	// Print again
	fmt.Print("InOrder traversal after removal: ")
	InOrder(root)
	fmt.Println()
}
