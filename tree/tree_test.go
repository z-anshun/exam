package tree

import (
	"exam/db"
	"testing"
)

func BenchmarkNode_IsContain(b *testing.B) {
	db.InitDb()
	s := db.GetWords()
	t := NewTree()
	for i := 0; i < len(s); i++ {
		//fmt.Println(s[i])
		t.AddNote(s[i])
	}

	k := "suicide"
	for i := 0; i < b.N; i++ {
		for _, v := range s {
			if v == k {
				break
			}
		}
	}

}
func BenchmarkNode_Tree(b *testing.B) {
	db.InitDb()
	s := db.GetWords()
	t := NewTree()
	for i := 0; i < len(s); i++ {
		//fmt.Println(s[i])
		t.AddNote(s[i])
	}
	k := "suicide"
	for i := 0; i < b.N; i++ {
		t.IsContain(k)
	}
}
