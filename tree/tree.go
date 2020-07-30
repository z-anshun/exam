package tree

type nodeType uint8

const (
	middle nodeType = iota //默认的中间节点
	root                   //根节点
	end                    //尾节点
)

var T = NewTree()        //根节点  这个是词汇数
var NameTree = NewTree() //这个是黑名单树

type node struct {
	//当前节点的值
	str string
	//指引 子节点的第一个
	indices []string
	//子节点
	children []*node
	nType    nodeType
	word     string
}

func NewTree() *node {
	return &node{
		nType: root,
	}
}

func (n *node) AddNote(str string) {
	if len(str) == 0 {
		return
	}
	l := len(n.str)
	//啥也没有
	if len(n.indices) == 0 && len(n.children) == 0 {

		child := &node{
			str:   str,
			nType: end,
			word:  str,
		}
		n.indices = append(n.indices, string(str[0]))
		n.children = append(n.children, child)
		return
	}

	//匹配的数
	i := longestCommonPrefix(n.str, str)

	//存在交集
	if i < len(n.str) {
		//这个是自己变后
		copy_n := &node{
			str:      n.str[i:],
			indices:  n.indices,
			children: n.children,
			nType:    n.nType,
			word:     n.word,
		}

		//直接塞
		n.children = []*node{copy_n}
		n.str = n.str[:i]
		n.indices = []string{string(copy_n.str[0])} //现加它自己的
		n.word = n.word[:(len(n.word) + i - l)]     //到目前的词汇
		n.nType = middle

		//塞新来的这个
		//如果path对于n.path有不同的
		if i < len(str) {
			str = str[i:]
			//进行塞
			n.addChild(str)
			return

		} else {
			//n.str 包含了 str
			n.nType = copy_n.nType

		}
		return
	} else { //str 包含n.str
		if i == len(str) {
			n.nType = end
			return
		}

		//遍历寻找
		for k := 0; k < len(n.indices); k++ {
			if n.indices[k] == string(str[i]) {
				n.children[k].AddNote(str[i:])
				return
			}
		}
		//没有匹配的
		n.addChild(str[i:])
	}

}

//是否包含
func (n *node) IsContain(str string) (string, bool) {
	for i := 0; i < len(str); i++ {
		word, b := n.FindStr(str[i:])
		if b {
			return word, b
		}
	}
	return "", false
}

//添加就只是添加就好
func (n *node) addChild(str string) {

	child := &node{
		str:   str,
		nType: end,
		word:  n.word + str,
	}
	n.children = append(n.children, child)
	n.indices = append(n.indices, string(str[0]))

}

//是否匹配
func (n *node) FindStr(str string) (string, bool) {

	switch n.nType {
	case root:
		return n.findChild(str)
	case middle:

		if len(str) < len(n.str) {
			return "", false
		}
		i := longestCommonPrefix(str, n.str)
		if i == len(n.str) && i < len(str) {

			return n.findChild(str[i:])
		}
		if i == len(n.str) && i == len(str) {
			return "", false
		}
	case end:

		if n.str == str || str == n.str[:len(n.str)-1] {

			return n.word, true
		}

		i := longestCommonPrefix(str, n.str)
		if len(n.str) < len(str) && i == len(str) && len(n.indices) != 0 {
			return n.findChild(str[i:])
		}

		return "", false
	}
	return "", false
}

func (n *node) findChild(str string) (string, bool) {
	for i := 0; i < len(n.indices); i++ {
		if n.indices[i] == string(str[0]) {
			return n.children[i].FindStr(str)
		}

	}

	return "", false
}

//最大匹配数
func longestCommonPrefix(nstr string, str string) int {
	longest := 0
	for ; longest < min(len(nstr), len(str)); longest++ {
		if nstr[longest] != str[longest] {
			return longest
		}
	}
	return longest
}

//获取较小的一个
func min(a int, b int) int {
	if a >= b {
		return b
	}
	return a
}
