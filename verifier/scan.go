//@Time  : 2018/5/22 22:35
//@Author: Greg Li
package verifier

import (
	"gitlab.com/usechain/go-usechain/crypto"
	"gitlab.com/usechain/go-usechain/accounts/keystore"
	"gitlab.com/usechain/go-usechain/common/hexutil"
	"fmt"
	"crypto/ecdsa"
)

// A1=[hash([b]A)]G+S
func GenerateSubAccount(bA *ecdsa.PublicKey, S *ecdsa.PublicKey) {
	hashBytes := crypto.Keccak256(crypto.FromECDSAPub(bA))   //hash([b]A)
	A1 := new(ecdsa.PublicKey)
	A1.Curve = crypto.S256()
	A1.X, A1.Y = crypto.S256().ScalarBaseMult(hashBytes)
	A1.X, A1.Y = crypto.S256().Add(A1.X, A1.Y, S.X, S.Y)
	fmt.Println(A1)

	address := crypto.PubkeyToAddress(*A1).Hex()
	fmt.Println(address)

	return
}


// A1=[hash([b]S)]G+A
var bPriv,_=hexutil.Decode("0x9d44c4709733e48669d37412a7257829bfaa3e74cfb575a408bf0b06de3b1d9a")

func ScanGenerateA1() ecdsa.PublicKey {
	A1S1:="0x03421ea288d37091f77d2ffce9180944913b94317804eb875b85797e03f3fff68002d5e258ff90035254d21549b66ba14680b03561a2bfe633548f18fb27d948aa04"
	sbyte,_:=hexutil.Decode(A1S1)
	A2, S1, err := keystore.GeneratePKPairFromABaddress(sbyte[:])
	if err !=nil {
		fmt.Println(err)
	}

	A1 := new(ecdsa.PublicKey)
	A1.X, A1.Y = crypto.S256().ScalarMult(S1.X, S1.Y, bPriv)   //A1=[b]S
	A1Bytes := crypto.Keccak256(crypto.FromECDSAPub(A1))        //hash([b]S)
	A1.X, A1.Y = crypto.S256().ScalarBaseMult(A1Bytes)   //[hash([b]S)]G

	AA,_:=hexutil.Decode("0x0402a99e51949742f41b37b7895b05b9dcaf75ed3e7e6952b39597cdb5f509cc43ec5a1157036e199d5ecbebbf2b6d07850525ae7d9f5f2d7de86fedfb69ec5380")
	A:=crypto.ToECDSAPub(AA)

	A1.X, A1.Y = crypto.S256().Add(A1.X, A1.Y, A.X, A.Y) //A1=[hash([b]S)]G+A
	A1.Curve = crypto.S256()

	////////////test A2变成公钥
	pubBytes := crypto.FromECDSAPub(A2)
	pub22:=hexutil.Encode(pubBytes)
	fmt.Println("pub22::::::::",pub22)
	address := crypto.PubkeyToAddress(*A2).Hex()
	fmt.Println(address)
	return *A1
}



func main()  {
	a:=ScanGenerateA1()
	fmt.Println(a)
	pubBytes := crypto.FromECDSAPub(&a)
	fmt.Println(hexutil.Encode(pubBytes))
	ABtoAddr()
}

func ABtoAddr()  {
	ab:="0x03ddc4b818f8c9991af603810f4be281efe0f1a737d60990421a02df647e8e11fa0370f7c1c284a7dc12769c1a18ddd3c56c7ee2cad1677c78c35b2ab74f5441f29a"
	abByte,_:=hexutil.Decode(ab)
	A, _, _ := keystore.GeneratePKPairFromABaddress(abByte)
	address := crypto.PubkeyToAddress(*A).Hex()
	fmt.Println(address)
}