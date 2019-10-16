#!/usr/bin/python3

"""Constructs the curl command to the geth RPC service.

Does not include the target URL.

"""

import os

transaction_hex = os.popen('ruby ./scripts/generate_tx.rb | grep -o 0x.* | tail -1'
).read().strip()

tx_cmd = f'''"method":"eth_sendRawTransaction","params":["{transaction_hex}"]'''
tx_data = f''' '{{"jsonrpc":"2.0",{tx_cmd},"id":1}}' '''
print(f'''curl -X POST -H "Content-Type: application/json" --data ''' +
      tx_data + " 172.16.1.100:8545")

# Sample output: curl -X POST -H "Content-Type: application/json" --data  '{"jsonrpc":"2.0","method":"eth_sendRawTransaction","params":[0xf8d2808504a817c800830186a08080b87e6080604052348015600f57600080fd5b507fa9f7ac292bec89f3379964cf3da3a5dec83186cd18c08abe2c0603a06fb19c2960405160405180910390a160358060496000396000f3006080604052600080fd00a165627a7a72305820753aa6316e04ef5d28155af1be96e95e873ef4ab750ebf088aaa47d11950e53f0029820fb7a0bbfd050ec4160068319d089f8ac8f288ea3f8fc11e56557060ef3cfbfb8701f9a06e3a3e045c9fd34de14bbfa11d631de873311a241efa1678efeade0513eb9725],"id":1}'
