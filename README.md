# optimism-tx-format

Shows how optimistic geth modifies transactions that get submitted to it by converting them to OVM Messages and then calls to the OVM Execution Manager


```bash
# it is assumed that optimism's geth is in the same parent dir as this repo to apply
# the go.mod replacement
git clone https://github.com/ethereum-optimism/go-ethereum
git clone https://github.com/gakonst/optimism-tx-format
cd optimism-tx-format
go run main.go
```

Expected output:

```
Message w/o OVM modding:
From: 0x71562b71999873DB5b286dF957af199Ec94617F7
GasPrice: 0
Value: 0
L1 Msg Sender: 0xaAaAaAaaAaAaAaaAaAAAAAAAAaaaAaAaAaaAaaAa
L1 Block number: 5
Queue Origin: 0
Nonce: 0
To: 0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB
Data: 0x0
Gas: 21000

OVM Message modded to pass through the Sequencer Entrypoint:
From: 0x71562b71999873DB5b286dF957af199Ec94617F7
GasPrice: 0
Value: 0
L1 Msg Sender: 0xaAaAaAaaAaAaAaaAaAAAAAAAAaaaAaAaAaaAaaAa
L1 Block number: 5
Queue Origin: 0
Nonce: 0
To: 0x1111111111111111111111111111111111111111
Data: 0x0053ee60121718bf75e0a93666613e1b9d099c5d0d2c2ef812ed2de0fe5e3e288371d3a9ff7049145104229577f9f474be9d1ef016940e407bd723a6dbacb749b401005208000000000000bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
Gas: 21000

OVM Message modded to pass through the Execution Manager:
From: 0x71562b71999873DB5b286dF957af199Ec94617F7
GasPrice: 0
Value: 0
L1 Msg Sender: 0xaAaAaAaaAaAaAaaAaAAAAAAAAaaaAaAaAaaAaaAa
L1 Block number: 5
Queue Origin: 0
Nonce: 0
To: 0x2222222222222222222222222222222222222222
Data: 0x9be3ad6700000000000000000000000000000000000000000000000000000000000000400000000000000000000000003333333333333333333333333333333333333333000000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000000000000000000000000000aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa0000000000000000000000001111111111111111111111111111111111111111000000000000000000000000000000000000000000000000000000000000520800000000000000000000000000000000000000000000000000000000000000e0000000000000000000000000000000000000000000000000000000000000005f0053ee60121718bf75e0a93666613e1b9d099c5d0d2c2ef812ed2de0fe5e3e288371d3a9ff7049145104229577f9f474be9d1ef016940e407bd723a6dbacb749b401005208000000000000bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb00
Gas: 125000000
```

As you can see, the `to`, `data` and `gas` fields of the Ovm Message get modded as it
propagates across the system, so that it passes first through the Sequencer Entrypoint
for calldata compression, and then through the Execution Manager for OVM sandboxing
