package fetch

import (
	"fmt"
	"log"
	"math/big"
	"sort"
	"sync"

	"github.com/0fatih/avalanche-balances/token"
	. "github.com/0fatih/avalanche-balances/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Balance struct {
	balance string
	name    string
}

var balances = []Balance{}

func FetchBalances(target string) {
	var wg sync.WaitGroup

	client, err := ethclient.Dial("wss://api.avax.network/ext/bc/C/ws")
	if err != nil {
		log.Fatal(err)
	}

	for _, tkn := range tokenAddresses {
		go getTokenBalance(&wg, client, tkn, target)
		wg.Add(1)
	}
	wg.Wait()

	sort.Slice(balances, func(i, j int) bool {
		return balances[i].balance > balances[j].balance
	})

	fmt.Printf("queried %d tokens\n\n", len(balances))
	for _, v := range balances {
		if v.balance == "0" {
			continue
		}
		balanceInEther := new(big.Int)
		balanceInEther, ok := balanceInEther.SetString(v.balance, 10)
		if !ok {
			fmt.Println("SetString: error")
			return
		}

		// Print beautiful
		fmt.Printf("%s: %s\n", v.name, FormatCommas(WeiToEther(balanceInEther)))
	}
}

func getTokenBalance(wg *sync.WaitGroup, client *ethclient.Client, tkn string, target string) {

	tokenAddress := common.HexToAddress(tkn)
	instance, err := token.NewToken(tokenAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	address := common.HexToAddress(target)
	bal, err := instance.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		log.Fatal(err)
	}

	symbol, err := instance.Symbol(&bind.CallOpts{})
	if err != nil {
		log.Fatal("Error in ", tokenAddress, err)
	}

	balances = append(balances, Balance{name: symbol, balance: bal.Text(10)})
	defer wg.Done()
}
