# OffchainLookup Notes

## OffchainConfig Logs

Below is logs showing the `OffchainConfig` being grabbed and used to add both headers and url param substitutions. First you see the offchain revert data. Notice the `url` has `{version}` which will get replaced. You can also see
the `fields` which are jq style field descriptions for json parsing.

```
[DEBUG] [OffchainLookup]: {url:https://pokeapi.co/api/{version}/pokemon/1 extraData:[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 32 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 4 48 120 48 48 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] fields:[.id .name [.abilities[] | .ability.name] [.types[] | .type.name]|join(",")] callbackFunction:[183 114 215 10]}
```

Then we see the cache hit for the `UpkeepInfo` which holds the `OffchainConfig`. To see the config you can take a like at what I added in the cli tool.

```
[DEBUG] [OffchainLookup] cache hit UpkeepInfo: {Target:0x57779C309fB2Bc4A87C5dbd181dD043A8aAC825b ExecuteGas:5000000 CheckData:[48 120 48 48] Balance:+10000000000000000000 Admin:0x00894C8b2B1a5635d2014287D2762a751Be0DcBe MaxValidBlocknumber:4294967295 LastPerformBlockNumber:0 AmountSpent:+0 Paused:false OffchainConfig:[161 100 107 101 121 115 130 164 100 116 121 112 101 102 72 101 97 100 101 114 100 110 97 109 101 112 88 45 84 101 115 116 105 110 103 45 72 101 97 100 101 114 101 118 97 108 117 101 88 51 130 88 32 81 59 6 155 103 29 163 187 143 246 253 114 7 181 164 185 4 112 85 213 17 251 215 56 16 100 131 143 170 237 189 6 111 84 72 73 83 95 73 83 95 65 80 73 95 75 69 89 106 68 101 99 114 121 112 116 86 97 108 111 84 72 73 83 95 73 83 95 65 80 73 95 75 69 89 164 100 116 121 112 101 101 112 97 114 97 109 100 110 97 109 101 103 118 101 114 115 105 111 110 101 118 97 108 117 101 88 38 130 88 32 81 59 6 155 103 29 163 187 143 246 253 114 7 181 164 185 4 112 85 213 17 251 215 56 16 100 131 143 170 237 189 6 98 118 50 106 68 101 99 114 121 112 116 86 97 108 98 118 50]}
```

```go
offchainConfig = OffchainAPIKeys{Keys: []Key{
    {
        Type:       "Header",
        Name:       "X-Testing-Header",
        Value:      nil,
        DecryptVal: "THIS_IS_API_KEY",
    },
    {
        Type:       "param",
        Name:       "version",
        Value:      nil,
        DecryptVal: "v2",
    },
}}
```

Next we see the headers added and also the url change as `{version}` was replaced with `v2`. And then we see the external request succeed and the json parsed for values. You can see the parsed values aims to give back the user whatever they request so the user could build a multi-value like `[.abilities[] | .ability.name]` which says "create an array of all abilities ability names" ie `[overgrow chlorophyll]` or `[.types[] | .type.name]|join(",")` which says "create an array of type names and join tha values with a comma-space" ie `grass,poison`. The user could format the response for further parsing in their contract.

```
[INFO]  [OffchainLookup] Headers: map[Content-Type:[application/json] X-Testing-Header:[THIS_IS_API_KEY]]
[INFO]  [OffchainLookup] URL: https://pokeapi.co/api/v2/pokemon/1
[DEBUG] [OffchainLookup] StatusCode: 200
[DEBUG] [OffchainLookup] Parsed values: [1 bulbasaur [overgrow chlorophyll] grass,poison]
```

Here is another uninterrupted example

```
[DEBUG] [OffchainLookup]: {url:https://pokeapi.co/api/{version}/pokemon/6 extraData:[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 32 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 4 48 120 48 48 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] fields:[.id .name [.abilities[] | .ability.name] [.types[] | .type.name]|join(",")] callbackFunction:[183 114 215 10]}
[DEBUG] [OffchainLookup] cache hit UpkeepInfo: {Target:0x57779C309fB2Bc4A87C5dbd181dD043A8aAC825b ExecuteGas:5000000 CheckData:[48 120 48 48] Balance:+10000000000000000000 Admin:0x00894C8b2B1a5635d2014287D2762a751Be0DcBe MaxValidBlocknumber:4294967295 LastPerformBlockNumber:0 AmountSpent:+0 Paused:false OffchainConfig:[161 100 107 101 121 115 130 164 100 116 121 112 101 102 72 101 97 100 101 114 100 110 97 109 101 112 88 45 84 101 115 116 105 110 103 45 72 101 97 100 101 114 101 118 97 108 117 101 88 51 130 88 32 81 59 6 155 103 29 163 187 143 246 253 114 7 181 164 185 4 112 85 213 17 251 215 56 16 100 131 143 170 237 189 6 111 84 72 73 83 95 73 83 95 65 80 73 95 75 69 89 106 68 101 99 114 121 112 116 86 97 108 111 84 72 73 83 95 73 83 95 65 80 73 95 75 69 89 164 100 116 121 112 101 101 112 97 114 97 109 100 110 97 109 101 103 118 101 114 115 105 111 110 101 118 97 108 117 101 88 38 130 88 32 81 59 6 155 103 29 163 187 143 246 253 114 7 181 164 185 4 112 85 213 17 251 215 56 16 100 131 143 170 237 189 6 98 118 50 106 68 101 99 114 121 112 116 86 97 108 98 118 50]}
[DEBUG] [OffchainLookup] Headers: map[Content-Type:[application/json] X-Testing-Header:[THIS_IS_API_KEY]]
[DEBUG] [OffchainLookup] URL: https://pokeapi.co/api/v2/pokemon/6
[DEBUG] [OffchainLookup] StatusCode: 200
[DEBUG] [OffchainLookup] Parsed values: [6 charizard [blaze solar-power] fire,flying]
[DEBUG] [OffchainLookup] Success: {Key:8638514|36741630312195675343567778757922113014347294099056117006653153098173984390406 State:1 FailureReason:0 GasUsed:+43189 PerformData:[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 128 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 192 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 64 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 54 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 9 99 104 97 114 105 122 97 114 100 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 19 91 98 108 97 122 101 32 115 111 108 97 114 45 112 111 119 101 114 93 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 11 102 105 114 101 44 102 108 121 105 110 103 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] FastGasWei:+2000000000 LinkNative:+4188778405487047 CheckBlockNumber:8638513 CheckBlockHash:[54 15 61 29 46 75 135 213 67 204 213 232 227 139 174 79 106 213 80 123 149 198 242 25 188 216 33 131 14 232 106 78] ExecuteGas:5000000}
```

## Show offchain config not ready

here we can see what happens when offchain config is not set yet and offchain lookup fails. I note this because since a user will need to register their upkeep and then  update offchainConfig this will happen if the upkeep is set to perform before that is updated if they need api keys else this is of no concern.

```
[DEBUG] [OffchainLookup]: {url:https://pokeapi.co/api/{version}/pokemon/1 extraData:[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 32 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 4 48 120 48 48 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] fields:[.id .name [.abilities[] | .ability.name] [.types[] | .type.name]|join(",")] callbackFunction:[183 114 215 10]}
[DEBUG] [OffchainLookup] cache hit UpkeepInfo: {Target:0x57779C309fB2Bc4A87C5dbd181dD043A8aAC825b ExecuteGas:5000000 CheckData:[48 120 48 48] Balance:+10000000000000000000 Admin:0x00894C8b2B1a5635d2014287D2762a751Be0DcBe MaxValidBlocknumber:4294967295 LastPerformBlockNumber:0 AmountSpent:+0 Paused:false OffchainConfig:[]}
[DEBUG] [OffchainLookup] Headers: map[Content-Type:[application/json]]
[DEBUG] [OffchainLookup] URL: https://pokeapi.co/api/%7Bversion%7D/pokemon/1
[DEBUG] [OffchainLookup] StatusCode: 404
```

since the param substitution was unable to happen the request was still attempted but the path is wrong so it 404s. If the param substitution is an actual url param like `http://example.com?apkikey={api-key}` then the path would work but
probably get a 401/403.

## Show events for upkeep

Here I show my test upkeep successful event logs showing the offchain lookup working

https://goerli.etherscan.io/address/0x57779c309fb2bc4a87c5dbd181dd043a8aac825b#events

I added a command to the cli to read the sample contract
```shell
❯ go run main.go keeper upkeep-pokemon-events 0x57779C309fB2Bc4A87C5dbd181dD043A8aAC825b 8638490 8638563 --config goerli.env
2023/03/11 14:03:49 Using config file goerli.env
```
note: there is only 4 parsed json fields but due to the last fields value being a comma separate list some rows look to have an extra field. So I added quotes for clarity 
```csv
from, id, name, abilities, types
0x0292f1C3219Ce8568dB1ee3b68090a66bf9e0151,"1","bulbasaur","[overgrow chlorophyll]","grass,poison"
0xBb4abA5Ea8D166f5eAab9Ac8a9eBA47E7BB95Db2,"2","ivysaur","[overgrow chlorophyll]","grass,poison"
0x0292f1C3219Ce8568dB1ee3b68090a66bf9e0151,"3","venusaur","[overgrow chlorophyll]","grass,poison"
0x02E2116F926668BCBA7c84BF824748D0107A3931,"4","charmander","[blaze solar-power]","fire"
0xBb4abA5Ea8D166f5eAab9Ac8a9eBA47E7BB95Db2,"5","charmeleon","[blaze solar-power]","fire"
0xBb4abA5Ea8D166f5eAab9Ac8a9eBA47E7BB95Db2,"6","charizard","[blaze solar-power]","fire,flying"
0xBb4abA5Ea8D166f5eAab9Ac8a9eBA47E7BB95Db2,"7","squirtle","[torrent rain-dish]","water"
0xBb4abA5Ea8D166f5eAab9Ac8a9eBA47E7BB95Db2,"8","wartortle","[torrent rain-dish]","water"
0xBb4abA5Ea8D166f5eAab9Ac8a9eBA47E7BB95Db2,"9","blastoise","[torrent rain-dish]","water"
0x4935A8c2a504272B9309Bde8cAD83F21d3646Ca5,"10","caterpie","[shield-dust run-away]","bug"
0xBb4abA5Ea8D166f5eAab9Ac8a9eBA47E7BB95Db2,"11","metapod","[shed-skin]","bug"
0x4935A8c2a504272B9309Bde8cAD83F21d3646Ca5,"12","butterfree","[compound-eyes tinted-lens]","bug,flying"
```
If you want to see 326 pokemon events
```shell
go run main.go keeper upkeep-pokemon-events 0x57779C309fB2Bc4A87C5dbd181dD043A8aAC825b 8638490 8640238 --config goerli.env
```
