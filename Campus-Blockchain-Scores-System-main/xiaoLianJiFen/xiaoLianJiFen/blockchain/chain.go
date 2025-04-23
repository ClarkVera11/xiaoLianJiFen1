package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"time"

	contracts "xiaoLianJiFen/go-contract" // 修改为基于模块名的绝对路径
	Models "xiaoLianJiFen/models"

	"github.com/astaxie/beego/orm"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 用你 Remix 部署后的地址替换
var contractAddress = common.HexToAddress("0xF7af7E9a21Abb0413112aadFbe2A734D86F1C302")

// UploadAllUsersToBlockchain 将所有用户信息上传到区块链
func UploadAllUsersToBlockchain(privateKey *ecdsa.PrivateKey) error {
	o := orm.NewOrm()
	var users []Models.Users

	// 获取所有用户
	_, err := o.QueryTable("users").All(&users)
	if err != nil {
		return fmt.Errorf("获取用户列表失败: %v", err)
	}

	log.Printf("开始上传 %d 个用户的信息到区块链", len(users))

	// 遍历所有用户并上传到区块链
	for _, user := range users {
		// 如果用户已经有交易哈希，说明已经上链，跳过
		if user.TxHash != "" {
			log.Printf("用户 %s 已经上链，跳过", user.Username)
			continue
		}

		err := UploadStudentPoints(privateKey, user.Username)
		if err != nil {
			log.Printf("用户 %s 上链失败: %v", user.Username, err)
			continue
		}

		log.Printf("用户 %s 成功上链", user.Username)
		// 等待一段时间，避免交易拥堵
		time.Sleep(time.Second * 2)
	}

	return nil
}

// 上链方法，传入用户信息
func UploadStudentPoints(privateKey *ecdsa.PrivateKey, username string) error {
	// 从数据库获取用户信息
	o := orm.NewOrm()
	var user Models.Users
	err := o.QueryTable("users").Filter("username", username).One(&user)
	if err != nil {
		return fmt.Errorf("获取用户信息失败: %v", err)
	}

	client, err := ethclient.Dial("http://127.0.0.1:8545") // 可替换为本地链
	if err != nil {
		return err
	}

	// 用你测试账户私钥替换 ↓↓↓

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return err
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	log.Println("Sender address:", fromAddress.Hex())
	balance, _ := client.BalanceAt(context.Background(), fromAddress, nil)
	log.Println("账户余额（wei）:", balance)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	chainID := big.NewInt(887766)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(3000000)
	auth.GasPrice = gasPrice

	instance, err := contracts.NewStudentPoints(contractAddress, client)
	if err != nil {
		return err
	}

	// 调用合约方法
	tx, err := instance.UpdateUserInfo(auth, user.Username, user.Phone, user.Role_name, big.NewInt(int64(user.Points)), user.Title)
	if err != nil {
		return err
	}

	log.Printf("交易已发送，等待确认，交易哈希: %s", tx.Hash().Hex())

	// 等待交易被确认，最多等待 10 次，每次 10 秒
	var receipt *types.Receipt
	for i := 0; i < 10; i++ {
		receipt, err = client.TransactionReceipt(context.Background(), tx.Hash())
		if err == nil && receipt != nil {
			// 交易已确认
			break
		}
		log.Printf("等待交易被确认...第 %d 次重试", i+1)
		time.Sleep(10 * time.Second)
	}

	if err != nil || receipt == nil {
		return fmt.Errorf("交易确认失败: %v", err)
	}

	if receipt.Status == 0 {
		return fmt.Errorf("交易执行失败")
	}

	// 获取区块时间戳
	block, err := client.BlockByNumber(context.Background(), receipt.BlockNumber)
	if err != nil {
		return fmt.Errorf("获取区块信息失败: %v", err)
	}

	// 更新数据库中的交易哈希和时间戳
	user.TxHash = tx.Hash().Hex()
	user.BlockTimestamp = int64(block.Time())
	_, err = o.Update(&user, "TxHash", "BlockTimestamp")
	if err != nil {
		return fmt.Errorf("更新数据库失败: %v", err)
	}

	log.Printf("✅ 交易已确认，数据库更新成功，交易哈希: %s, 时间戳: %d", tx.Hash().Hex(), block.Time())
	return nil
}
