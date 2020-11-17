var fs = require('fs');
const toml = require('toml');
var json2toml = require('json2toml');
var data = JSON.parse(fs.readFileSync('keys.json', 'utf8'));

const spec = {
    "AlphaPPB": 10000000,
    "DeltaC": "120s",
    "DeltaGrace": "10s",
    "DeltaProgress": "45s",
    "DeltaStage": "30s",
    "DeltaResend": "10s",
    "DeltaRound": "20s",
    "RMax": 4,
    "OracleIdentities": [
        
    ]
}
// makeConfig reads the keys.json file in local directory, loops through all nodes except for the first node(bootstrap node) and build the data.config file for the prototype tool.
function makeConfig(){
    var data = JSON.parse(fs.readFileSync('keys.json', 'utf8'));
    for(let x in data){        
        if(x==1){
        }else{
        var obj = 
            {
                "TransmitAddress": data[x]["ETH"][0].address, 
                "OnchainSigningAddress": data[x]["OCR"][0].onChainSigningAddress, 
                "PeerID": data[x]["P2P"][0].peerId, 
                "OffchainPublicKey": data[x]["OCR"][0].offChainPublicKey,
                "ConfigPublicKey": data[x]["OCR"][0].configPublicKey
            };
        spec["OracleIdentities"].push(obj);     
        }
    }
    console.log(spec)
    fs.writeFile("config/data.config", JSON.stringify(spec, null, 4), function (err) {
        if (err) throw err;
        console.log('Created data.config!');
      });
}
//makeConfig()

function createToml(){
    const jobSpec = toml.parse(fs.readFileSync('jobspec_2.toml', 'utf-8'));
    var dataConverted = json2toml(jobSpec);
    fs.writeFile("jobs/jobspec_2.toml", dataConverted, function (err) {
    if (err) throw err;
    });
}
function createJobspecs(contractAddress){
    if(contractAddress == undefined){
        throw new ReferenceError("Must provide a contractAddress.")
    }
    var data = JSON.parse(fs.readFileSync('keys.json', 'utf8'));
    const bootstrapJob = toml.parse(fs.readFileSync('bootstrap.toml', 'utf-8'));    
    const nodeJob = toml.parse(fs.readFileSync('node.toml', 'utf-8'));
    for(x=1;x<6; x++){        
        if(x==1){            
            bootstrapJob["contractAddress"]=contractAddress;
            bootstrapJob["p2pPeerID"]=data[1]["P2P"][0].peerId;
            var dataConverted = json2toml(bootstrapJob);
            fs.writeFile("jobs/jobspec_1.toml", dataConverted, function (err) {
            if (err) throw err;
            });
        }else{
            nodeJob["contractAddress"]=contractAddress;
            nodeJob["p2pPeerID"]=data[x]["P2P"][0].peerId;
            nodeJob["p2pBootstrapPeers"]=["/dns4/chainlink-node-1/tcp/6690/p2p/"+data[1]["P2P"][0].peerId];
            nodeJob["keyBundleID"]=data[x]["OCR"][0].id;
            nodeJob["transmitterAddress"]=data[x]["ETH"][0].address;
            var dataConverted = json2toml(nodeJob);
            fs.writeFile("jobs/jobspec_"+x+".toml", dataConverted, function (err) {
            if (err) throw err;
            });
        }
    }
    console.log("Jobs Done!")
}
createJobspecs("0xD4988baC598e872F326fBEDab223E7B8cA3CEe95")