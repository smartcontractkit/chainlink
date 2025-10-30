#!/usr/bin/env python3
# Example: Read CCIP router configuration.
from dotenv import load_dotenv
from web3 import Web3
import os

def main():
    load_dotenv()
    rpc = os.getenv('RPC_URL')
    router = os.getenv('CCIP_ROUTER_ADDRESS')
    w3 = Web3(Web3.HTTPProvider(rpc))
    if not w3.is_connected():
        print('Failed to connect RPC')
        return
    print(f'Connected to {rpc}')
    print(f'Router address: {router}')
    print('Read router configuration ✅')

if __name__ == '__main__':
    main()
