package middlewares

import (
	"fmt"

	"github.com/Codexiaoyi/linweb/interfaces"
)

func Logs(context interfaces.IContext) {
	fmt.Println("Request path: ", context.Request().Path())
	//处理请求
	context.Next()
}
