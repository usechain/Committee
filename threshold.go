package sssa

import (
	"fmt"
	"gitlab.com/usechain/go-usechain/crypto"
	"crypto/ecdsa"
	"math/big"
	"gitlab.com/usechain/go-usechain/common"
	"gitlab.com/usechain/go-usechain/common/hexutil"
	"gitlab.com/usechain/go-usechain/commitee/verifier"

)

var Sharespart []*big.Int = make([]*big.Int, 3)

func generatePrivKey(key *big.Int) *ecdsa.PrivateKey {
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = crypto.S256()
	priv.D = key.Mod(key, crypto.S256().Params().N)
	priv.PublicKey.X, priv.PublicKey.Y = crypto.S256().ScalarBaseMult(key.Bytes())

	return priv
}

//count sA
func multiPub(s []byte, A *ecdsa.PublicKey) *ecdsa.PublicKey {
	A1 := new(ecdsa.PublicKey)
	A1.Curve = crypto.S256()
	//以A为基点加s次
	A1.X, A1.Y = crypto.S256().ScalarMult(A.X, A.Y, s)   //A1=[s]B
	return A1
}

func GenerateShares(minimum int, shares int) ([]ecdsa.PublicKey,[]*big.Int) {
	// Short, medium, and long tests
	strings := "N17FigASkL6p1EOgJhRaIquQLGvY"

	created, pointer, polynomial, err := Create(minimum, shares, strings)
	if err != nil {
		fmt.Println("Fatal: creating: ", err)
	}

	var pubArray []ecdsa.PublicKey
	for k := range polynomial {
		fmt.Printf( "The polynomial: %x,%d\n", polynomial[k],k)
		priv := generatePrivKey(polynomial[k])
		pubArray = append(pubArray, priv.PublicKey)
		fmt.Printf("The public key is :%x\n", pubArray[k])
	}

	pubSum := new(ecdsa.PublicKey)
	pubSum.Curve = crypto.S256()

	for k := range polynomial {
		if k == 0 {
			pubSum.X = pubArray[0].X
			pubSum.Y = pubArray[0].Y
			continue
		}
		pubSum.X, pubSum.Y = crypto.S256().Add(pubSum.X, pubSum.Y, pubArray[k].X, pubArray[k].Y)
	}
	fmt.Printf("The sum key is:%x\n\n\n", pubSum)


	for j := range pointer {
		//fmt.Println("The created num:", created[j],)

		priv := generatePrivKey(pointer[j])
		pubKey := priv.PublicKey
		fmt.Printf("The pointer is %d, %x, :%x\n", j+1, pointer[j], pubKey)
	}
	combined, err := Combine(created)
	fmt.Println("The combined num:", combined)
	if err != nil {
		fmt.Println("Fatal: combining: ", err)
	}


	return pubArray, pointer
}

func GenerateSubAccountShares(serverId uint16) *ecdsa.PublicKey {
	sharePriv := ReadSelfshares(serverId)
	mainAccountPub, _ := GetMainAccountPub()
	shareAccountPub := crypto.ToECDSAPub(mainAccountPub)
	shareData := CountSubAccountSharePart(sharePriv, shareAccountPub)
	return shareData
}

func expontUint16(x uint16, y int) int {
	z := x
	for i := 1; i < y; i++ {
		z = z * x
	}
	return int(z)
}

func Checkshares(polynomial []PolynomialMsg, pointYstr string,  serverId uint16, senderId uint16) {
	var polynomialInt []ecdsa.PublicKey = make([]ecdsa.PublicKey, len(polynomial))

	fmt.Println("YYYYYYYYYYYYYYYYYYYY")
	for i := range polynomial {
		polynomialInt[i].Curve = crypto.S256()
		polynomialInt[i].X = fromBase64(polynomial[i].X)
		polynomialInt[i].Y = fromBase64(polynomial[i].Y)
		fmt.Printf("%x, %x\n", polynomialInt[i].X, polynomialInt[i].Y)
	}
	fmt.Printf("%x\n", fromBase64(pointYstr))
	fmt.Println("YYYYYYYYYYYYYYYYYYYY")

	pubSum := new(ecdsa.PublicKey)
	pubSum.Curve = crypto.S256()
	for j := range polynomialInt {
		fmt.Println(polynomialInt[j].X, polynomialInt[j].Y)
		if j == 0 {
			pubSum.X = polynomialInt[j].X
			pubSum.Y = polynomialInt[j].Y
			continue
		}
		for k := 0; k < expontUint16(serverId, j); k++ {
			pubSum.X, pubSum.Y = crypto.S256().Add(pubSum.X, pubSum.Y, polynomialInt[j].X, polynomialInt[j].Y)
		}

	}
	fmt.Printf("The sum key is:%x\n\n\n", pubSum)

	priv := generatePrivKey(fromBase64(pointYstr))
	pubKey := priv.PublicKey
	fmt.Printf("The pubkey is %x\n", pubKey)

	if pubSum.X.Cmp(pubKey.X) == 0 && pubSum.Y.Cmp(pubKey.Y) == 0 {
		fmt.Println("The shares is legal!")
	}

	// ...add it to results...
	fmt.Println("the id is", senderId-1)
	Sharespart[senderId-1] = fromBase64(pointYstr)
	return
}

func HandleSubAccountVerifyRequest(polynomial []PolynomialMsg, serverId uint16, senderId uint16, serverPort uint16) {
	var polynomialInt []ecdsa.PublicKey = make([]ecdsa.PublicKey, len(polynomial))

	for i := range polynomial {
		polynomialInt[i].Curve = crypto.S256()
		polynomialInt[i].X = fromBase64(polynomial[i].X)
		polynomialInt[i].Y = fromBase64(polynomial[i].Y)
		fmt.Printf("%x, %x\n", polynomialInt[i].X, polynomialInt[i].Y)
	}

	if serverId == 1 {
		fmt.Println("Get total pub:", polynomialInt[0])
		verifier.GenerateSubAccount(&polynomialInt[0], &polynomialInt[0])
	}else {
		pubSum := GenerateSubAccountShares(serverId)
		for j := range polynomialInt {
			fmt.Println(polynomialInt[j].X, polynomialInt[j].Y)
			if j == 0 {
				pubSum.X = polynomialInt[j].X
				pubSum.Y = polynomialInt[j].Y
				continue
			}
			for k := 0; k < expontUint16(serverId, j); k++ {
				pubSum.X, pubSum.Y = crypto.S256().Add(pubSum.X, pubSum.Y, polynomialInt[j].X, polynomialInt[j].Y)
			}

		}
		fmt.Printf("The sum key is:%x\n\n\n", pubSum)
		fmt.Println("Server:", serverPort, serverId)
		destID, destPoint := GetDestNode(serverPort, serverId)
		fmt.Println("Dest:", destID, destPoint)
		SendVerifyMsg(destPoint, destID, pubSum)
	}

	return
}


func CountSharesPart(id uint16) {
	shareSum := big.NewInt(0)

	for i := range Sharespart {
		shareSum.Add(shareSum, Sharespart[i])
	}

	// ...add it to results...
	result := ToBase64(big.NewInt(int64(id)))
	result += ToBase64(shareSum)
	fmt.Println("The shares base64: ", result)
}

func ReadSelfshares(id uint16) *big.Int {
	shares := []string{
		"3d1f7b376c2a9a58fe9ee622f4bc1886a3d8be4eee409ebf3d5b54053295f705",
		"cf5cf4a9f57a618709dff5369b17357cbcfa7b165c995ce2c536c86aa78b89fa",
		"8dd2076aa1fc541f23ea6da2123ca565c50783c0ae1c3b064d123cd01c811dac",
	}
	shareInt, _ := big.NewInt(0).SetString(shares[id - 1], 16)
	return shareInt
}

func GetMainAccountPub() ([]byte, error) {
	Astr := "0x049a0b2c928af39a0dd635702e920864d16ec9846d1517a5e181792d4b84943688746359d46c49045d42b550a27f464919c1838f93d478750deeec48a8a9db12a6"
	return hexutil.Decode(Astr)
}


func CountSubAccountSharePart(sharePriv *big.Int,key *ecdsa.PublicKey)  *ecdsa.PublicKey {
	return multiPub(sharePriv.Bytes(), key)
}


func EestLibraryCombine() {
	shares := []string{
		//"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=PR97N2wqmlj-nuYi9LwYhqPYvk7uQJ6_PVtUBTKV9wU=",
		"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE=z1z0qfV6YYcJ3_U2mxc1fLz6exZcmVzixTbIaqeLifo=",
		"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAI=jdIHaqH8VB8j6m2iEjylZcUHg8CuHDsGTRI80ByBHaw=",
	}

	combined, err := Combine(shares)
	if err != nil {
		fmt.Println("Fatal: combining: ", err)
	}
	fmt.Printf("The combined string: %x\n", combined)

	hexStr := fmt.Sprintf("%x", combined)

	privInt,_ := big.NewInt(0).SetString(hexStr, 16)

	priv := generatePrivKey(privInt)
	pubKey := priv.PublicKey
	fmt.Printf("The pubkey is %x\n", pubKey)

	pub:=common.ToHex(crypto.FromECDSAPub(&pubKey))
	fmt.Println(pub)

	if combined != "test-pass" {
		fmt.Println("Failed library cross-language check")
	}
}




