package clear
import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)
// 创建一个命令行清除函数变量
var cl map[string]func() 
func init(){
	//初始化map
	cl = make(map[string]func())
	cl["linux"] = func ()  {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	cl["windows"] = func ()  {
		cmd := exec.Command("cmd","/c","cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear(){
	value,ok := cl[runtime.GOOS]
	if ok{
		value()
	}else{
		fmt.Println("该平台不支持清除命令行输出>_<!")
	}
}

func ClearCmd(){
	CallClear()
}






