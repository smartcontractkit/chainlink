# Make Money, Hit Me On My Beeper

[![Hit Me On My Beeper](HitMeOnMyBeeper.jpg?raw=true 'Hit Me On My Beeper')](https://www.youtube.com/watch?v=-4Wu-zSndlw&t=13s)

In this example, we will send a text message when a contract address receives ether.
This will link Ethereum events (on-chain) with the [Twilio](https://www.twilio.com) SMS service
(off-chain) using Chainlink (CL)

## Configure and run [Chainlink development environment](../README.md#run-chainlink-development-environment)

## Sign up for a free Twilio account

- https://www.twilio.com/try-twilio

## Run Node Server

This Node JS server relays messages from Chainlink jobs to Twilio.

- `yarn install`
- `node twilio.js <twilio_account_sid> <twilio_auth_token> <twilio_number> <your_number>`
  - i.e. `./twilio.js AC97fade171cxxxxxxxxxxxxxxxxxxxxxx ffab4f5ecc65acxxxxxxefe99xxxxx "+1 786-555-5555" 3055555555`.

## Create Chainlink Job and Smart Contract (where you will send the money)

- `truffle migrate`, remember the contract address (0x...) of the migration above.
- `./create_twilio_job_for 0x...` pass contract address as argument

## Send Money

- `./send_money_to 0x...` pass contract address as argument

## Celebrate

[![Hit Me On My Beeper](HitMeOnMyBeeper.jpg?raw=true 'Hit Me On My Beeper')](https://www.youtube.com/watch?v=-4Wu-zSndlw&t=13s)
