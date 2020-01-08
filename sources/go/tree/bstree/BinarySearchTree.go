package bstree

import (
	"fmt"
)

// 节点对象
type node struct {
	data	int
	left	*node
	right 	*node
	parent	*node				// 记录父节点的原因是可以简化删除操作，以及对红黑树的演化
}

// 树对象
type BSTree struct {
	root	*node
	length 	int
}

func NewBSTree() *BSTree{
	return &BSTree{
		root:   nil,
		length: 0,
	}
}

// 插入元素：迭代方式
func (t *BSTree)Insert(data int) {

	// 如果根节点为空
	if t.root == nil {
		root := &node{
			data:   data,
			left:   nil,
			right:  nil,
			parent: nil,
		}
		t.root = root
		t.length++
		return
	}

	// 如果添加的不是根节点
	currentNode := t.root			// 循环到的当前节点
	parentNode := t.root			// 当前节点的父节点
	isLeft := true					// 当前是否是在左侧插入

	// 创建插入节点
	insertNode := &node{
		data:   data,
		left:   nil,
		right:  nil,
		parent: nil,
	}

	// 找到要插入的位置
	for currentNode != nil {			// 插入其实是直接往叶节点处插入，所以循环退出条件是到达叶节点后

		// 插入节点与当前节点的数据相等
		if insertNode.data == currentNode.data {
			currentNode = insertNode	// 万一以后要将data属性改变为自定义对象呢？所以这里要做一下相等处理
			return
		}

		// 插入节点与当前节点的数据不相等
		parentNode = currentNode
		if insertNode.data > currentNode.data {		//  往右插入
			currentNode = currentNode.right
			isLeft = false
		} else {									// 往左插入
			currentNode = currentNode.left
			isLeft = true
		}
	}

	// 执行插入
	insertNode.parent = parentNode
	if isLeft {
		parentNode.left = insertNode
	} else {
		parentNode.right = insertNode
	}
	t.length++
}

// 打印二叉搜索树：从上往下打印，其实质是层序遍历
func (t *BSTree)Display() {
	levelOrderTraverse(t.root)
}

// 查找节点：根据值查找值所在的节点
func (t *BSTree)findNode(data int) *node {

	currentNode := t.root
	for currentNode != nil {
		if data == currentNode.data  {
			return currentNode
		}
		if data > currentNode.data {
			currentNode = currentNode.right
		} else {
			currentNode = currentNode.left
		}
	}

	return currentNode
}

// 删除元素
func (t *BSTree)Remove(data int) bool {

	// 查找元素值对应的节点
	delNode := t.findNode(data)
	if delNode == nil {
		fmt.Println("未找到该节点")
		return false
	}

	if delNode.left != nil && delNode.right != nil {				// 删除节点有2个子节点
		t.removeTwoNode(delNode)
	} else if delNode.left == nil && delNode.right == nil{			// 删除节点有0个子节点
		t.removeZeroNode(delNode)
	} else {														// 删除节点有1个子节点
		t.removeOneNode(delNode)
	}

	t.length--
	return true
}

// 删除叶节点
func (t *BSTree)removeZeroNode(n *node){

	if n.left != nil || n.right != nil {
		panic("传入非法节点")
	}

	if n == t.root {
		t.root = nil
	}

	if n.parent.left == n {
		n.parent.left = nil
	} else {
		n.parent.right = nil
	}
}

// 删除度为1节点
func (t *BSTree)removeOneNode(n *node){

	if (n.left == nil && n.right == nil) || (n.left != nil && n.right != nil){
		panic("传入非法节点")
	}

	var replace *node
	if n.left != nil {
		replace = n.left
	} else {
		replace = n.right
	}

	replace.parent = n.parent

	// 如果 n 是根节点
	if n == t.root {
		t.root = replace
		return
	}

	if n.parent.left == n {
		n.parent.left = replace
	} else {
		n.parent.right = replace
	}
}

// 删除度为2节点
func (t *BSTree)removeTwoNode(n *node){

	if n.left == nil || n.right == nil {
		panic("传入非法节点")
	}

	// 查找前驱(后继也可以)
	prev := t.findPrevNode(n)
	if prev == nil {
		panic("未找到后继节点")
	}

	fmt.Println("找到的前驱节点：", prev.data)

	// 删除当前节点
	n.data = prev.data
	n = prev

	// 如果当前节点是根节点
	if n == t.root {
		n.parent = nil
	}

	// 删除前驱节点:前驱节点必定是度为1或0的节点
	n.parent = prev.parent
	if prev.left != nil || prev.right != nil{
		t.removeOneNode(prev)
	} else {
		t.removeZeroNode(prev)
	}
}

/*
因为二叉搜索树获取最大值、最小很方便，所以需要提供最大值、最小值获取操作
*/
// 获取最大值：二叉搜索树最右边的值
func (t *BSTree)MaxData() interface{}{
	currentNode := t.root
	var data interface{}
	for currentNode != nil {
		data = currentNode.data
		currentNode = currentNode.right
	}
	return data
}

// 获取最小值：二叉搜索树最左边的值
func (t *BSTree)MinData() interface{}{
	currentNode := t.root
	var data interface{}
	for currentNode != nil {
		data = currentNode.data
		currentNode = currentNode.left
	}
	return data
}

// 查找数据是否在二叉树
func (t *BSTree)SearchData(data int) bool {
	currentNode := t.root
	for currentNode != nil {
		if data < currentNode.data {
			currentNode = currentNode.left
		} else if data > currentNode.data{
			currentNode = currentNode.right
		} else {
			return true
		}
	}
	return false
}


/**
	节点查找相关方法：适合所有二叉树
 */

// 查找前驱节点：即二叉树在中序遍历时，当前元素的前一个节点，二叉搜索树的前驱节点也是当前节点前的最大节点
func (t *BSTree)findPrevNode(n *node) *node{

	// 情况一：当前节点为空
	if n == nil {
		fmt.Println("传入节点为空")
		return nil
	}

	// 情况二：当前节点的左子节点为空，父节点也为空
	if n.left == nil && n.parent == nil {
		fmt.Println("无前驱节点")
		return nil
	}

	// 定义当前循环到了哪个节点
	var currentNode *node

	// 情况三：当前节点的左子节点为空，父节点不为空，前驱在在parent的右子树中
	if n.left == nil && n.parent != nil {
		currentNode = n
		for currentNode.parent != nil && currentNode == currentNode.parent.left {
			currentNode = currentNode.parent
		}
	} else {			// 最后一个情况： n.left != nil
		currentNode = n.left
		for currentNode.left != nil {
			currentNode = currentNode.right
		}
	}

	return currentNode
}

// 查找后继节点：即二叉树在中序遍历时，当前元素的后一个节点，二叉搜索树的后继节点也是当前节点后的最小节点
func (t *BSTree)findNextNode(n *node) *node {

	// 情况一：当前节点为空
	if n == nil {
		fmt.Println("传入节点为空")
		return nil
	}

	// 情况二：当前节点的左子节点为空，父节点也为空
	if n.right == nil && n.parent == nil {
		fmt.Println("无前驱节点")
		return nil
	}

	// 定义当前循环到了哪个节点
	var currentNode *node

	// 情况三：当前节点的左子节点为空，父节点不为空，前驱在在parent的右子树中
	if n.right == nil && n.parent != nil {
		currentNode = n
		for currentNode.parent != nil && currentNode == currentNode.parent.right {
			currentNode = currentNode.parent
		}
	} else {			// 最后一个情况： n.right != nil
		currentNode = n.right
		for currentNode.right != nil {
			currentNode = currentNode.left
		}
	}

	return currentNode
}


// 插入元素：递归方式
func (t *BSTree)InsertByRC(data int) {

	// 构造插入节点
	insertNode := &node{
		data:  data,
		left:  nil,
		right: nil,
		parent: nil,
	}

	// 判断根节点是否存在
	if t.root == nil {
		t.root = insertNode
	} else {
		// 执行递归插入
		insertRC(t.root, insertNode)
	}

	t.length++
}
