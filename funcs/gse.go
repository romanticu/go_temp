package funcs

import (
	"fmt"

	"github.com/go-ego/gse"
)

func SegTest() {
	// text := "GPRS系统中CCU与PCU间的数据传输实现方法"
	var seg gse.Segmenter
	seg.LoadDict()
	// 分词文本
	tb := "门体联锁开关 采用该门体联锁开关 微波烹饪设"
	fmt.Println("输出分词结果, 类型为字符串, 使用搜索模式: ", seg.String(tb, true))
	fmt.Println("输出分词结果, 类型为 slice: ", seg.Slice(tb))

	segments := seg.Segment([]byte(tb))
	// 处理分词结果, 普通模式
	
	fmt.Println(gse.ToString(segments))
	fmt.Println(gse.ToSlice(segments))
	po := seg.Pos(tb)
	fmt.Println("pos: ", po)

}
