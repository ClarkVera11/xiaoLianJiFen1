package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"xiaoLianJiFen/blockchain"
	_ "xiaoLianJiFen/models"
	_ "xiaoLianJiFen/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	_ "github.com/go-sql-driver/mysql"
)

func add1(a int) int {
	return a + 1
}

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	//orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/xiaoLianJiFen?charset=utf8")
}

func main() {

	// 读取 keystore 文件内容
	files, err := os.ReadDir(`D:/geth student 1009/student/dev-chain/keystore`)
	if err != nil {
		log.Fatal("读取目录失败:", err)
	}
	keystoreFile := `D:/geth student 1009/student/dev-chain/keystore/` + files[0].Name()
	password := "12345678" // 输入你的钱包密码

	// 打开 keystore 文件
	file, err := os.Open(keystoreFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 解密 keystore 文件
	keyjson, err := os.ReadFile(keystoreFile)
	if err != nil {
		log.Fatal(err)
	}

	// 解密为私钥
	key, err := keystore.DecryptKey(keyjson, password)
	if err != nil {
		log.Fatal(err)
	}

	// 获取私钥
	privateKey := key.PrivateKey

	// 输出私钥
	fmt.Printf("Private Key: %x\n", crypto.FromECDSA(privateKey))

	err = beego.AddFuncMap("ShowPrePage", HandlePerPage)
	if err != nil {

	}
	err = beego.AddFuncMap("ShowNextPage", HandleNextPage)
	if err != nil {
		return
	}

	// 上传所有用户信息到区块链
	err = blockchain.UploadAllUsersToBlockchain(privateKey)
	if err != nil {
		beego.Error("批量上链失败：", err)
	} else {
		beego.Info("所有用户信息成功上链")
	}

	beego.AddFuncMap("add1", add1)
	beego.Run()
}

func HandlePerPage(data int) string {

	pageIndex := data - 1

	pageIndex1 := strconv.Itoa(pageIndex)
	return pageIndex1
}

func HandleNextPage(data int) string {

	pageIndex := data + 1

	pageIndex1 := strconv.Itoa(pageIndex)
	return pageIndex1
}
