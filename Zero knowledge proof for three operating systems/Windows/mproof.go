package main

import (
	"fmt"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc"
	bn254 "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
)

type Circuit struct {
	PreImage frontend.Variable //Hkey
	Ht       frontend.Variable `gnark:",public"`
}

func (circuit *Circuit) Define(curveID ecc.ID, api frontend.API) error {
	Htt, _ := mimc.NewMiMC(TKeyt, curveID, api)
	Htt.Write(circuit.PreImage)
	Httb := api.ToBinary(Htt.Sum())
	Htb := api.ToBinary(circuit.Ht)
	for i := 0; i < len(Httb); i++ {
		api.AssertIsEqual(Httb[i], Htb[i])
	}
	return nil
}

var Key = "544893542849cbded60ce6880112235267e0908f580e9db65e138663fb15f78f" //sha256
var t = "11111111"
var T = "00000000"
var TKeyt = Key + t + T
var HASH = "16725690776755845349664889100637449378411911100218711130668876130116619486512"

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
	respvk := Makepvk()
	if respvk {
		fmt.Println("Zero knowledge proof initialization successful")
	}
	proof, respro := Makeproof()
	Proof = proof
	if respro {
		fmt.Println("Zero knowledge proof generated successfully")
	}
	resver := Makeverify()
	if resver {
		fmt.Println("Verification successful")
	} else {
		fmt.Println("Verification failed")
	}
}

func Makepvk() bool {
	var circuit Circuit
	r1cs, err := frontend.Compile(ecc.BN254, backend.GROTH16, &circuit)
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
	preImage := []byte{0x01, 0x02, 0x03}
	hash := mimcHash(preImage, Key)
	fmt.Printf("hash(key): %s\n", hash)
	if hash != HASH {
		return nil, false
	} else {
		preImageHt := []byte{0x01, 0x02, 0x03}
		ht := mimcHash(preImageHt, TKeyt)
		fmt.Printf("Ht: %s\n", ht)

		witness := &Circuit{
			PreImage: frontend.Value(preImage),
			Ht:       frontend.Value(ht),
		}
		proof, err := groth16.Prove(R1cs, Pk, witness)
		if err != nil {
			fmt.Printf("Prove failed： %v\n", err)
			return nil, false
		}
		Proof = proof
		return proof, true
	}
}
func Makeverify() bool {
	preImageHt := []byte{0x01, 0x02, 0x03}
	ht := mimcHash(preImageHt, TKeyt)

	publicWitness := &Circuit{
		Ht: frontend.Value(ht),
	}

	err := groth16.Verify(Proof, Vk, publicWitness)
	if err != nil {
		fmt.Printf("verification failed: %v\n", err)
		return false
	}
	return true
}