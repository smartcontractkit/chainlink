#!/usr/bin/env python3
# Example: Monitor CCIP transaction
import time

def main():
    print('Monitoring CCIP transaction...')
    for i in range(3):
        time.sleep(1)
        print(f'Polling status... {i+1}/3')
    print('Transaction confirmed ✅')

if __name__ == '__main__':
    main()
