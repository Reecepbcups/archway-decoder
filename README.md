# arch-decoder

This just reduces overhead from the normal app so it ONLY decodes transactions. Used to speed up <https://github.com/Reecepbcups/cosmos-indexer>

```bash
make install

# Binaries have this, but this reduces overhead
# You can find these txs with https://rpc.network.com/block?height=7793647
arch-decode tx decode [amino-base64]

# Example in action
arch-decode tx decode --output json CrUBCrIBCiQvY29zbXdhc20ud2FzbS52MS5Nc2dFeGVjdXRlQ29udHJhY3QSiQEKK2p1bm8xamc5ZjNkbnJzZzNkN25wNTRyanlxY3J6MHphZXcydjNyejVlNGQSP2p1bm8xNjRjZDZ3cGVwa3hwNG5kd2t2d2U5bnFwOWtyYzM3bHduNzB0N3h5M2g3amFhNHlwZWtzc3p1MnczdxoZeyJkaXN0cmlidXRlX3Jld2FyZHMiOnt9fRJpClIKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiECJoHdaIdtGcDLsWxbEngOKq6GU6S3ykOXO8HeXt5H6EMSBAoCCAEYvtgDEhMKDQoFdWp1bm8SBDE2MDIQteMJGkAM8dn7/9tazy1CJGClB/GLP5TXzJIxUTcKHnP93kQKBGKBDqkxEUUbzzVtoNYLw/K0CMSzsMKibVQjVK6D7xfR
```

---

## Mass decode file
Useful if you have a lot of data and want multithreaded performance. Input and output uses unique key ids to ensure you get the correct data

**arch-decode tx decode-file input.json out.json**
Then you can read

Input:
```json
[
    {
        "id":1,
        "tx": "ClMKUQobL2Nvc21vcy5nb3YudjFiZXRhMS5Nc2dWb3RlEjIInAISK2p1bm8xNmR6bjRwd3Q4cjZ3cm42ODc4OGNrY2g5ajdrMnF6eWxmdXVlOXkYARJmClEKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiED1Vfp+F7xleWzjZElWfoubOAhFCcZy0Ocg25CVx0O2kISBAoCCAEY/wUSEQoLCgV1anVubxICNzUQ8sgEGkBqQboZOTd+1Yai6SkuRskq+LadkgSnlKY/YdHEizoYAW0HGboRQQFhILFsWdJVlOPeNYdIP/QE9/n9cJUjortp"
    },
    {
        "id":1,
        "tx": "CrsBCrgBCiQvY29zbXdhc20ud2FzbS52MS5Nc2dFeGVjdXRlQ29udHJhY3QSjwEKK2p1bm8xeWF2OTY4eTZtenJnN25semR3OTVnYTRjeWR6cmR5dmpoYTN0NTgSP2p1bm8xOHd1eTVxcjJtc3dnejd6YWs4eXI5Y3Jod2h0dXIzdjZtdzR0Y3l0dXB5d3h6dzdzdWZ5cWd6YTd1aBoLeyJib25kIjp7fX0qEgoFdWp1bm8SCTEwMDAwMDAwMBJpClEKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiEDvAPkqUCXD7uf0apB/6vON1cwFPEFPsiiPno08RPsOpcSBAoCCH8Y4gESFAoOCgV1anVubxIFNzA5OTkQ0+M5GkCbf9otpKEi0n1W4buJ2ne+k7KPdbtca5aUSza1lWERFVCw3bQcUUS6FanDsh43Ez0np14VT/vUV3FIft0ExqDB"
    }
]
```

Output:
(use json.loads() in python to parse said Txs)
```json
[
    {
        "id": 1,
        "tx": "{\"body\":{\"messages\":[{\"@type\":\"/cosmwasm.wasm.v1.MsgExecuteContract\",\"sender\":\"juno1yav968y6mzrg7nlzdw95ga4cydzrdyvjha3t58\",\"contract\":\"juno18wuy5qr2mswgz7zak8yr9crhwhtur3v6mw4tcytupywxzw7sufyqgza7uh\",\"msg\":{\"bond\":{}},\"funds\":[{\"denom\":\"ujuno\",\"amount\":\"100000000\"}]}],\"memo\":\"\",\"timeout_height\":\"0\",\"extension_options\":[],\"non_critical_extension_options\":[]},\"auth_info\":{\"signer_infos\":[{\"public_key\":{\"@type\":\"/cosmos.crypto.secp256k1.PubKey\",\"key\":\"A7wD5KlAlw+7n9GqQf+rzjdXMBTxBT7Ioj56NPET7DqX\"},\"mode_info\":{\"single\":{\"mode\":\"SIGN_MODE_LEGACY_AMINO_JSON\"}},\"sequence\":\"226\"}],\"fee\":{\"amount\":[{\"denom\":\"ujuno\",\"amount\":\"70999\"}],\"gas_limit\":\"946643\",\"payer\":\"\",\"granter\":\"\"}},\"signatures\":[\"m3/aLaShItJ9VuG7idp3vpOyj3W7XGuWlEs2tZVhERVQsN20HFFEuhWpw7IeNxM9J6deFU/71FdxSH7dBMagwQ==\"]}"
    },
    {
        "id": 1,
        "tx": "{\"body\":{\"messages\":[{\"@type\":\"/cosmos.gov.v1beta1.MsgVote\",\"proposal_id\":\"284\",\"voter\":\"juno16dzn4pwt8r6wrn68788ckch9j7k2qzylfuue9y\",\"option\":\"VOTE_OPTION_YES\"}],\"memo\":\"\",\"timeout_height\":\"0\",\"extension_options\":[],\"non_critical_extension_options\":[]},\"auth_info\":{\"signer_infos\":[{\"public_key\":{\"@type\":\"/cosmos.crypto.secp256k1.PubKey\",\"key\":\"A9VX6fhe8ZXls42RJVn6LmzgIRQnGctDnINuQlcdDtpC\"},\"mode_info\":{\"single\":{\"mode\":\"SIGN_MODE_DIRECT\"}},\"sequence\":\"767\"}],\"fee\":{\"amount\":[{\"denom\":\"ujuno\",\"amount\":\"75\"}],\"gas_limit\":\"74866\",\"payer\":\"\",\"granter\":\"\"}},\"signatures\":[\"akG6GTk3ftWGoukpLkbJKvi2nZIEp5SmP2HRxIs6GAFtBxm6EUEBYSCxbFnSVZTj3jWHSD/0BPf5/XCVI6K7aQ==\"]}"
    }
]
```
