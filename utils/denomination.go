package utils

import "math/big"

const Ether = 1000000000000000000 // 1 Ether in wei

func EtherToWei(val *big.Int) *big.Int {
	return new(big.Int).Mul(val, big.NewInt(Ether))
}

func WeiToEther(val *big.Int) *big.Int {
	return new(big.Int).Div(val, big.NewInt(Ether))
}
