
package main

import (
	"fmt"
	"math/big"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	bn254 "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
)

type Circuit struct {
	PreImageSalt frontend.Variable //T
	PreImage     frontend.Variable //key
	Timestamp    frontend.Variable `gnark:",public"`
	Ht           frontend.Variable `gnark:",public"`
	HASH         frontend.Variable `gnark:",public"`
}

func (circuit *Circuit) Define(curveID ecc.ID, api frontend.API) error {
	sum := api.Add(circuit.PreImageSalt, circuit.PreImage, circuit.Timestamp)
	mimcHash, _ := mimc.NewMiMC("seed", curveID, api)
	// Step 3: 对加法结果 sum 进行哈希
	mimcHash.Write(sum)
	// Step 4: 获取哈希结果并与公开的哈希值进行比较
	Hk, _ := mimc.NewMiMC("seed", curveID, api)
	Hk.Write(circuit.PreImage)
	hash2result := Hk.Sum()
	ee := api.IsZero(api.Sub(hash2result, circuit.HASH))
	hashResult := api.Select(ee, mimcHash.Sum(), 0)
	api.AssertIsEqual(hashResult, circuit.Ht)

	// Hk, _ := mimc.NewMiMC("seed", curveID, api)
	// Hk.Write(circuit.PreImage)
	// hash2result := Hk.Sum()
	// Select(hash2result==circuit.HASH, 0, mimcHash.Sum())
	// hash2resultb := api.FromBinary(api.Xor(api.ToBinary(hash2result), api.ToBinary(HASH)))

	// api.AssertIsEqual(0, hash2resultb)

	return nil
}

type NullifierRegistry struct {
	nullifierSet map[string]bool
}

// NewNullifierRegistry 创建一个新的 NullifierRegistry 实例
func NewNullifierRegistry() *NullifierRegistry {
	return &NullifierRegistry{
		nullifierSet: make(map[string]bool),
	}
}

// AddNullifier 添加新的 Nullifier
// 如果已经存在，返回 false 表示重复提交，否则返回 true 表示成功添加
func (nr *NullifierRegistry) AddNullifier(nullifier string) bool {
	if nr.nullifierSet[nullifier] {
		return false // Nullifier 已存在，不能重复使用
	}
	nr.nullifierSet[nullifier] = true
	return true // 成功添加新的 Nullifier
}

// IsNullifierUnique 检查 Nullifier 是否已经提交过
func (nr *NullifierRegistry) IsNullifierUnique(nullifier string) bool {
	return !nr.nullifierSet[nullifier]
}

var Key = "38122509095249768280721192201792540884047770874958887343840368830908906993551" //"544893542849cbded60ce6880112235267e0908f580e9db65e138663fb15f78f" 的十进制形式
var t = "1728975932271510400"
var T = "00000000"

// HASH is the hash value of Key, which can be calculated by the following codes
// HASH := mimcHash(bKey.Bytes(), "seed")
// fmt.Printf("hash(key): %s\n", hash)
var HASH = "5827313906892128376317691618591248229023467324282166300794902454604441795955"

// Ht is the sum of Key,t, and T , which can be calculated by the following codes
// tInt := new(big.Int)
// TInt := new(big.Int)
// tInt, _ = tInt.SetString(t, 10)
// TInt, _ = TInt.SetString(T, 10)
// sum := new(big.Int).Add(bKey, tInt)
// sum = sum.Add(sum, TInt)
// ht := mimcHash(sum.Bytes(), "seed")
// fmt.Printf("Ht: %s\n", ht)
var Ht = "19727640260386863455988554025984213913104425721451513328001445945877966072857"
// 创建一个表示 6666 的 big.Int
var bigExpireTime = big.NewInt(11497821885134) //expireTime
func mimcHash(data []byte, str string) string {
	f := bn254.NewMiMC(str)
	f.Write(data)
	hash := f.Sum(nil)
	hashInt := big.NewInt(0).SetBytes(hash)
	return hashInt.String()
}

var R1cs frontend.CompiledConstraintSystem
var Pk groth16.ProvingKey
var Vk groth16.VerifyingKey
var Proof groth16.Proof

func main() {
	tInt := new(big.Int)
	TInt := new(big.Int)
	tInt, _ = tInt.SetString(t, 10)
	TInt, _ = TInt.SetString(T, 10)
	bKey1, _ := new(big.Int).SetString(Key, 10)
	sum1 := new(big.Int).Add(bKey1, tInt)
	sum2 := new(big.Int).Add(sum1, TInt)
	Hash1 := mimcHash(bKey1.Bytes(), "seed")
	Hash2 := mimcHash(sum2.Bytes(), "seed")
	fmt.Printf("hash): %s\n", Hash1)
	fmt.Printf("Ht): %s\n", Hash2)
	respvk := Makepvk()
	if respvk {
		fmt.Println("Zero knowledge proof initialization successful")
	}
	startTime := time.Now()
	for i := 0; i < 100; i++ {
		_, _ = Makeproof()
	}
	proof, respro := Makeproof()
	// 记录执行后的时间戳
	endTime := time.Now()

	// 计算时间差，并转换为毫秒
	elapsedTime := float64(endTime.Sub(startTime).Nanoseconds()) / 1e6
	fmt.Printf("函数 Makeproof() 执行耗时: %.2f ms\n", elapsedTime)
	if !respro {
		fmt.Println("Zero knowledge proof is not generated successfully")
		return
	}
	Proof = proof
	if respro {
		fmt.Println("Zero knowledge proof generated successfully")
	}
	startTime2 := time.Now()
	for i := 0; i < 100; i++ {
		resver := Makeverify()
		if resver {
			fmt.Println("Verification successful")
		} else {
			fmt.Println("Verification failed")
		}
	}
	endTime2 := time.Now()
	elapsedTime2 := float64(endTime2.Sub(startTime2).Nanoseconds()) / 1e6
	fmt.Printf("函数 Makeverify() 执行耗时: %.4f ms\n", elapsedTime2)

}

func Makepvk() bool {
	var circuit Circuit
	r1cs, err := frontend.Compile(ecc.BN254, backend.GROTH16, &circuit)
	fmt.Printf("The size of circuit is : %d \n", r1cs.GetNbConstraints())
	if err != nil {
		fmt.Printf("Compile failed : %v\n", err)
		return false
	}

	pk, vk, err := groth16.Setup(r1cs)
	if err != nil {
		fmt.Printf("Setup failed\n")
		return false
	}
	R1cs = r1cs
	Pk = pk
	Vk = vk
	return true
}
func Makeproof() (groth16.Proof, bool) {
	bKey, _ := new(big.Int).SetString(Key, 10)

	tInt := new(big.Int)
	TInt := new(big.Int)

	tInt, _ = tInt.SetString(t, 10)
	TInt, _ = TInt.SetString(T, 10)

	bhash, _ := new(big.Int).SetString(HASH, 10)
	bHt, _ := new(big.Int).SetString(Ht, 10)
	witness := &Circuit{
		PreImageSalt: frontend.Value(TInt),
		Ht:           frontend.Value(bHt),
		PreImage:     frontend.Value(bKey),
		HASH:         frontend.Value(bhash),
		Timestamp:    frontend.Value(tInt),
	}
	proof, err := groth16.Prove(R1cs, Pk, witness)
	if err != nil {
		fmt.Printf("Prove failed： %v\n", err)
		return nil, false
	}
	Proof = proof
	return proof, true

}
func Makeverify() bool {
	registry := NewNullifierRegistry()
	tInt := new(big.Int)
	tInt, _ = tInt.SetString(t, 10)
	bHt, _ := new(big.Int).SetString(Ht, 10)
	if bHt.Int64() == 0 {
		return false
	}
	if registry.IsNullifierUnique(Ht) {
		registry.AddNullifier(Ht)
	} else {
		fmt.Printf("replay fail \n")
		return false
	}
	now := time.Now()
	unixNano := now.UnixNano()
	currenttime := new(big.Int).SetInt64(unixNano)
	fmt.Printf("currenttime:%s\n", currenttime)
	result := new(big.Int).Sub(currenttime, tInt)

	

	// 比较 result 和 6666
	comparison := result.Cmp(bigExpireTime)
	if comparison == 1 {
		fmt.Printf("Current time error : %s\n", new(big.Int).Sub(result, big6666))
		//return false
	}
	bhash, _ := new(big.Int).SetString(HASH, 10)
	publicWitness := &Circuit{
		Ht:        frontend.Value(bHt),
		HASH:      frontend.Value(bhash),
		Timestamp: frontend.Value(tInt),
	}

	err := groth16.Verify(Proof, Vk, publicWitness)
	if err != nil {
		fmt.Printf("verification failed: %v\n", err)
		return false
	}
	return true
}
