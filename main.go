package main

import (
	"fmt"
	"math/big"
	"os"

	"github.com/gakonst/optimism-tx-format/ovm"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
)

var testKey, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")

var testAddr = common.HexToAddress("0xBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB")
var testAddr2 = common.HexToAddress("0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")

// precompiles
var sequencerEntrypoint = common.HexToAddress("0x1111111111111111111111111111111111111111")
var executionManager = common.HexToAddress("0x2222222222222222222222222222222222222222")
var stateManager = common.HexToAddress("0x3333333333333333333333333333333333333333")

func debug(msg ovm.Message) {
	// these remain the same
	fmt.Printf("From: %+v\n", msg.From().Hex())
	fmt.Printf("GasPrice: %+v\n", msg.GasPrice().Int64())
	fmt.Printf("Value: %+v\n", msg.Value().Int64())
	fmt.Printf("L1 Msg Sender: %+v\n", msg.L1MessageSender().Hex())
	fmt.Printf("L1 Block number: %+v\n", msg.L1BlockNumber().Int64())
	fmt.Printf("Queue Origin: %+v\n", msg.QueueOrigin().Int64())
	fmt.Printf("Nonce: %+v\n", msg.Nonce())

	// these changes due to the mods
	fmt.Printf("To: %+v\n", msg.To().Hex())
	fmt.Printf("Data: %+v\n", common.ToHex(msg.Data()))
	fmt.Printf("Gas: %+v\n", msg.Gas())

	fmt.Println()
}

func main() {
	emAbiJson, err := os.Open("./em.json")
	if err != nil {
		panic(err)
	}
	defer emAbiJson.Close()
	emAbi, err := abi.JSON(emAbiJson)

	tx := types.NewTransaction(
		// normal Ethereum args
		0,
		testAddr,
		big.NewInt(0),
		21000,
		big.NewInt(0),
		[]byte{},

		// == Optimism args ==

		// L1 Message Sender
		&testAddr2,
		// L1 Block Number
		big.NewInt(5),
		// Queue origin (sequencer or L1 queue):
		// Transactions which come from the sequencer get modded.
		// If this were `types.QueueOriginL1ToL2`, the OVM message would not
		// be different.
		types.QueueOriginSequencer,
		// Sighash mode can either be EIP155 or EthSign. EIP155 is the "usual"
		// one used by all signing on Ethereum today, but EthSign is also added
		// so that transactions can also be signed via metamask's sign APIs.
		types.SighashEIP155,
	)

	// simple mainnet signer
	signer := types.NewEIP155Signer(big.NewInt(1))
	tx, err = types.SignTx(tx, signer, testKey)
	if err != nil {
		panic(err)
	}

	// The "normal" message will have no data field and the receiver will be the
	// one we specified
	msg, err := tx.AsMessage(signer)
	if err != nil {
		panic(err)
	}
	fmt.Println("Message w/o OVM modding:")
	debug(msg)

	// The OVM message will have a populated data field with the compressed
	// calldata which will be later published to Ethereum L1, and its `to` field
	// is rewired to be sent to the Sequencer Entrypoint
	ovmMsg, err := ovm.AsOvmMessage(tx, signer, sequencerEntrypoint)
	if err != nil {
		panic(err)
	}
	fmt.Println("OVM Message modded to pass through the Sequencer Entrypoint:")
	debug(ovmMsg)

	// Finally, when the transaction gets through the Execution Manager it will be
	// further modded

	// we mock the EVM's exec context
	evm := vm.EVM{}
	evm.Context.Time = big.NewInt(3)
	evm.Context.BlockNumber = big.NewInt(2)
	evm.Context.GasLimit = 125000000

	evm.Context.OvmExecutionManager.ABI = emAbi
	evm.Context.OvmExecutionManager.Address = executionManager
	evm.Context.OvmStateManager.Address = stateManager

	emMsg, err := ovm.ToExecutionManagerRun(&evm, ovmMsg)
	if err != nil {
		panic(err)
	}
	fmt.Println("OVM Message modded to pass through the Execution Manager:")
	debug(emMsg)
}
