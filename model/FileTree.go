package model

// FileTreeNode 如果类名首字母大写，表示其他包也能够访问
type FileTreeNode struct {
	parentPath    string
	name          string
	parentNode    *FileTreeNode
	childNode     []FileTreeNode
	baiduFileMap  map[string]File
	nativeFileMap map[string]File
}

func newFileTreeNode(parentPath string, name string, parentNode *FileTreeNode) FileTreeNode {
	var childNode []FileTreeNode
	return FileTreeNode{parentPath: parentPath, name: name, parentNode: parentNode, childNode: childNode}
}

func (receiver FileTreeNode) addNode(parentPath string, name string, parentNode *FileTreeNode) {
	fileNode := newFileTreeNode(parentPath, name, parentNode)
	receiver.childNode = append(receiver.childNode, fileNode)
}
