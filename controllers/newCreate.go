package controllers

import (
	"fmt"

	"github.com/emanueldf/form3-accounts/utils"
)

func Init() {
	fmt.Println(utils.ReadPublicKey())
	fmt.Println(utils.ReadPrivateKey())
}
