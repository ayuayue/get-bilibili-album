package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//Resp 定义响应信息结构体
type Resp struct {
	Code    int
	Msg     string
	Message string
	Data    Data
}

// Data 数据
type Data struct {
	Items []Item
}

// Item Data子集
type Item struct {
	DocId       int
	PosterUid   int
	Title       string
	Description string
	Pictures    []Picture
	Count       int
	Ctime       int
	View        int
	Like        int
}

// Picture info
type Picture struct {
	ImgSrc    string `json:"img_src"`
	ImgWidth  int    `json:"img_width"`
	ImgHeight int    `json:"img_height"`
	ImgSize   int    `json:"img_size"`
}

// Info uid
type Info struct {
	Uid      int
	PageNum  int
	PageSize int
	Biz      string
}

// Num 图片数量
type Num struct {
	Code    int
	Msg     string
	Message string
	Data    NumData
}

// NumData 图片信息
type NumData struct {
	AllCount   int `json:"all_count"`
	DrawCount  int `json:"draw_count"`
	PhotoCount int `json:"photo_count"`
	DailyCount int `json:"daily_count"`
}

func main() {
	user := new(Info)

	for {
		fmt.Println("请输入b站up主的uid : ")
		_, err := fmt.Scanln(&user.Uid)
		if err != nil {
			fmt.Println(err)
		}
		num := GetImgNum(user.Uid)
		// fmt.Println(num)
		if num != 0 {
			GetSrc(user.Uid, num)
			break
		}
		fmt.Println("--------------------------")
		fmt.Println("该up主还没有上传相册哦 >_<~~")
		time.Sleep(time.Second * 2)

	}

}

//GetSrc 获取图片链接
func GetSrc(uid int, num int) {
	index := 1
	pageNum := "0"
	biz := "all"
	//判断文件夹是否存在
	ok, err := IsExists(strconv.Itoa(uid))
	if err != nil {
		panic("文件夹检测失败 >_<!")
	}
	if !ok {
		err := os.Mkdir(strconv.Itoa(uid), os.ModePerm)
		if err != nil {
			panic("创建文件夹失败")
		}
	}

	url := "https://api.vc.bilibili.com/link_draw/v1/doc/doc_list?uid=" + strconv.Itoa(uid) + "&page_num=" + pageNum + "&page_size=" + strconv.Itoa(num) + "&biz=" + biz
	resp, err := http.Get(url)
	if err != nil {
		panic("好像解析错误了呢 >_<　！")
	}
	res := new(Resp)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic("数据读取失败 >_<!")
	}
	if err := json.Unmarshal(body, &res); err == nil {
		//获取所有items
		var items []Item
		items = append(items, res.Data.Items...)
		// fmt.Println(items)
		for _, v := range items {
			//获取每一个picture
			fmt.Println("正在读取第" + strconv.Itoa(index) + "个相册,共" + strconv.Itoa(num) + "个相册")

			for k, v1 := range v.Pictures {

				//开始写入图片
				//获取图片的后缀,并更改名字重新保存
				dot := strings.LastIndex(v1.ImgSrc, ".")
				ext := v1.ImgSrc[dot:]
				//检测文件是否存在,如果存在则跳过,不存在则写入
				path := strconv.Itoa(uid) + "/" + strconv.Itoa(index) + "-" + strconv.Itoa(k+1) + ext
				ok, err := IsExists(path)
				if err != nil {
					panic("文件检测错误 >_<!")
				}
				if !ok {
					//读取图片原始二进制数据
					bin, err := http.Get(v1.ImgSrc)
					if err != nil {
						panic("访问图片失败 >_<!")
					}
					defer bin.Body.Close()
					imgBin, err := ioutil.ReadAll(bin.Body)
					if err != nil {
						panic("读取图片数据失败 >_<!")
					}
					f, err := os.OpenFile(strconv.Itoa(uid)+"/"+strconv.Itoa(index)+"-"+strconv.Itoa(k+1)+ext, os.O_CREATE, 0666)
					defer f.Close()
					if err != nil {
						panic("文件读取失败 >_<!")
					}
					_, err = f.Write([]byte(imgBin))
					if err != nil {
						panic("文件写入失败 >_<!")
					}
					fmt.Println("正在下载第" + strconv.Itoa(index) + "个相册的第" + strconv.Itoa(k+1) + "个图片,共" + strconv.Itoa(len(v.Pictures)) + "个图片")
				}
			}
			index++
		}

	} else {
		panic("相册详情解析失败了呢 >_<!")
	}

}

//GetImgNum 获取相册数量
func GetImgNum(uid int) int {
	url := "https://api.vc.bilibili.com/link_draw/v1/doc/upload_count?uid=" + strconv.Itoa(uid)
	num := new(Num)
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("好像解析错误了呢 >_< !")
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("数据映射失败 >_< !")
	}
	err = json.Unmarshal(body, &num)
	if err != nil {
		fmt.Println("相册数量解析出现错误了 >_< !")

	}
	imgNum := num.Data.AllCount
	fmt.Println("-----------------获取相册数量 start------------------")
	if imgNum > 300 {
		fmt.Println("改up主的相册有点多,请稍等一下哦 >_<~~")
	} else {
		fmt.Println("-----------------获取相册数量 ing--------------------")
		fmt.Println("-----------------获取相册 end-----------------------")
		fmt.Printf("共获取到该up主 %v 个相册 >_<~~\n", imgNum)
	}
	return imgNum

}

//IsExists 检测文件是否存在
func IsExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
