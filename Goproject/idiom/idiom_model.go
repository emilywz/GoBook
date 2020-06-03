package idiom

import "fmt"

//成语
type Idiom struct {
	Title      string
	Spell      string
	Content    string
	Sample     string
	Derivation string
}

//打印成语数据
func PrintIdiom(idiom Idiom) {
	if idiom.Title != "" {
		fmt.Printf("Title:%s\n", idiom.Title)
		fmt.Printf("Spell:%s\n", idiom.Spell)
		fmt.Printf("Sample:%s\n", idiom.Sample)
		fmt.Printf("Derivation:%s\n", idiom.Derivation)
		fmt.Printf("Content:%s\n", idiom.Content)
	} else {
		fmt.Println("未找到成语！")
	}
}

