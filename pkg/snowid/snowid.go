package snowid

import (
	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

const nodeID int64 = 1

// 想用第三方库直接去看官方的代码示例，别在网上东找西找  这东西用起来还是比较简单的

func Init() (err error) {
	node, err = snowflake.NewNode(nodeID)
	return
}

// GenID 生成ID时会上锁，确保不重复
func GenID() string {
	return node.Generate().String()
}
