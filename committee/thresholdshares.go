package main

import (
	"fmt"
	"gitlab.com/usechain/go-usechain/commitee"
	"gitlab.com/usechain/go-usechain/crypto"
	"crypto/ecdsa"
	"math/big"
)

func generatePrivKey(key *big.Int) *ecdsa.PrivateKey {
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = crypto.S256()
	priv.D = key.Mod(key, crypto.S256().Params().N)
	priv.PublicKey.X, priv.PublicKey.Y = crypto.S256().ScalarBaseMult(key.Bytes())

	return priv
}

func main() {

	// Short, medium, and long tests
	strings := []string{
		"N17FigASkL6p1EOgJhRaIquQLGvYV0",
		"0y10VAfmyH7GLQY6QccCSLKJi8iFgpcS",
	}

	minimum := []int{2,2}
	shares := []int{3,3}

	var pubshares [3][2]*big.Int

	for i := range strings {
		created, pointer, polynomial, err := sssa.Create(minimum[i], shares[i], strings[i])
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


		for j := range created {
			//fmt.Println("The created num:", created[j],)

			priv := generatePrivKey(pointer[j])
			pubKey := priv.PublicKey
			fmt.Printf("The pointer is %d, %x, :%x\n", j+1, pointer[j], pubKey)
			pubshares[j][i] = pointer[j]
		}
		combined, err := sssa.Combine(created)
		fmt.Println("The combined num:", combined)
		if err != nil {
			fmt.Println("Fatal: combining: ", err)
		}
		if combined != strings[i] {
			fmt.Println("Fatal: combining returned invalid data")
		}
	}


	var sharesPart []*big.Int = make([]*big.Int, 3)
	for i := range pubshares {
		//for j:= range pubshares[i] {
		//	fmt.Printf("The pointer is %d, %d, :%x\n", i, j, pubshares[i][j])
		//}

		sharesPart[i] = pubshares[i][0]
		sharesPart[i] = sharesPart[i].Add(sharesPart[i],pubshares[i][1])
		fmt.Printf("The shares hex: %d, %x\n", i+1, sharesPart[i])

		// ...add it to results...
		result := sssa.ToBase64(big.NewInt(int64(i)))
		result += sssa.ToBase64(sharesPart[i])
		fmt.Println("The shares base64: ", result)
	}

	//TestLibraryCombine(sharesPart)
	EestLibraryCombine()
}


func TestLibraryCombine(shares []*big.Int) {
	combined, err := sssa.Tombine(shares)
	if err != nil {
		fmt.Println("Fatal: combining: ", err)
	}
	fmt.Printf("The combined :%x\n", combined)
	if combined != "test-pass" {
		fmt.Printf("Failed library cross-language check")
	}
}


func EestLibraryCombine() {
	shares := []string{
		"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=Dx1sBPeBj6_khAgJ24VatRTD6dpc12WOZf-nQF4D7Yw=",
		//"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE=EtBK92EnEp0IgcPhWlUC01s88bEwS9qBFn9u4P5ncSA=",
		"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAI=CbT1F4xghFF900Mb1X1zKcPEGuflJ_9zxv82gZ7K9XE=",
	}

	combined, err := sssa.Combine(shares)
	if err != nil {
		fmt.Println("Fatal: combining: ", err)
	}
	fmt.Printf("The combined string: %x\n", combined)
	if combined != "test-pass" {
		fmt.Println("Failed library cross-language check")
	}
}




