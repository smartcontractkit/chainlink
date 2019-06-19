require 'eth'

Eth.configure do |config|
  config.chain_id = 1994
end

gwei = 10**9
contract_data = "0x6080604052348015600f57600080fd5b507fa9f7ac292bec89f3379964cf3da3a5dec83186cd18c08abe2c0603a06fb19c2960405160405180910390a160358060496000396000f3006080604052600080fd00a165627a7a72305820753aa6316e04ef5d28155af1be96e95e873ef4ab750ebf088aaa47d11950e53f0029"

key = Eth::Key.new priv: '34d2ee6c703f755f9a205e322c68b8ff3425d915072ca7483190ac69684e548c'
tx = Eth::Tx.new({
                   nonce: 0,
                   gas_price: (20 * gwei),
                   gas_limit: 100_000,
                   value: ARGV[0].to_i,
                   data: contract_data
                 })

tx.sign key

puts tx.to_h

puts "hex for #{tx.nonce}: #{tx.hex}"
