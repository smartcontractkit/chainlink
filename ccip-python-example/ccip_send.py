#!/usr/bin/env python3
# Simplified CCIP send script for example purposes.
from dotenv import load_dotenv
from web3 import Web3
import os

def main():
    load_dotenv()
    rpc = os.getenv('RPC_URL')
    router = os.getenv('CCIP_ROUTER_ADDRESS')
    print(f'Connecting to RPC {rpc}')
    w3 = Web3(Web3.HTTPProvider(rpc))
    if not w3.is_connected():
        print('RPC connection failed.')
        return
    print(f'Router: {router}')
    print('Pretend to send CCIP message here... ✅')

if __name__ == '__main__':
    main()
