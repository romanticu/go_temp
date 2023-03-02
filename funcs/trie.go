package funcs

import (
	"fmt"
)

type Trie struct {
	Children map[rune]*Trie
	IsEnd    bool
	FullStr  string
}

// 插入候选词，并构造 Trie
func (t *Trie) Insert(words []string) {
	for _, word := range words {
		node := t
		// chArr := strings.Split(word, "")

		for _, ch := range word {
			if node.Children == nil {
				node.Children = make(map[rune]*Trie)
			}
			if node.Children[ch] == nil {
				node.Children[ch] = &Trie{}
			}
			node = node.Children[ch]
		}
		node.IsEnd = true
		node.FullStr = word
	}

}

// 搜索字符前缀
func (t *Trie) SearchPrefix(prefix string) *Trie {
	node := t
	// chArr := strings.Split(prefix, "")
	for _, ch := range prefix {

		if node.Children[ch] == nil {
			return nil
		}
		node = node.Children[ch]
	}
	return node
}

func (t *Trie) Search(word string) []string {
	node := t.SearchPrefix(word)
	return t.PrintWords(node, []string{})
}

// 深度遍历
func (t *Trie) PrintWords(node *Trie, words []string) []string {
	if node == nil {
		return words
	}
	if node.IsEnd {
		words = append(words, node.FullStr)
	}
	for _, n := range node.Children {
		words = t.PrintWords(n, words)
	}

	return words
}

func TestTrie() {
	trie := Trie{}
	fmt.Println("候选词输入：")
	fmt.Println(`"人世间", "人生在世", "人世间电视剧全集免费观看", "人世间剧情介绍", "还是人吗", "人世间演员表"`)
	trie.Insert([]string{"人世间", "人生在世", "人世间电视剧全集免费观看", "人世间剧情介绍", "还是人吗", "人世间演员表"})

	fmt.Println("搜索：人世间 \n 搜索结果: ")
	fmt.Println(trie.Search("人世间"))
}
