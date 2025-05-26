package config

import "crypto/ecdsa"

// 导出变量必须以大写开头：PrivateKey 而不是 privateKey
var PrivateKey *ecdsa.PrivateKey
