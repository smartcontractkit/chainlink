const Api = provider => {
  const Nat = require("./nat");
  const Map = require("./map");
  const Bytes = require("./bytes");
  const keccak256s = require("./hash").keccak256s;
  const send = method => (...params) =>
    new Promise((resolve,reject) =>
      provider.send(method, params, (err, result) =>
        err ? reject(err) : resolve(result)));

  // TODO check inputs
  // TODO move to proper file
  const encodeABI = (type, value) => {
    if (type === "bytes") {
      const length = Bytes.length(value);
      const nextMul32 = (((length - 1) / 32 | 0) + 1) * 32;
      const lengthEncoded = encodeABI("uint256", Nat.fromNumber(length)).data;
      const bytesEncoded = Bytes.padRight(nextMul32, value);
      return {data: Bytes.concat(lengthEncoded, bytesEncoded), dynamic: true};
    } else if (type === "uint256" || type === "bytes32") {
      return {data: Bytes.pad(32, value), dynamic: false};
    } else {
      throw "Eth-lib can't encode ABI type " + type + " yet.";
    }
  }

  const sendTransaction = send("eth_sendTransaction");
  const sendRawTransaction = send("eth_sendRawTransaction");
  const getTransactionReceipt = send("eth_getTransactionReceipt");
  const compileSolidity = send("eth_compileSolidity");
  const call = send("eth_call");
  const getBalance = send("eth_getBalance");
  const accounts = send("eth_accounts");

  const removeEmptyTo = tx =>
    tx.to === "" || tx.to === "0x" ? Map.remove("to")(tx) : tx;

  const waitTransactionReceipt = getTransactionReceipt; // TODO: implement correctly

  const addTransactionDefaults = tx =>
    // Get basic defaults
    Promise.all([
      tx.chainId || send("net_version")(),
      tx.gasPrice || send("eth_gasPrice")(),
      tx.nonce || send("eth_getTransactionCount")(tx.from,"latest"),
      tx.value || "0x0",
      tx.data || "0x"])
    // Add them to tx
    .then(([chainId, gasPrice, nonce, value, data]) => {
      return Map.merge(tx)({chainId: Nat.fromNumber(chainId), gasPrice, nonce, value, data});
    })
    .then(tx => {
      // Add gas default by estimating
      if (!tx.gas) {
        let estimateTx = {};
        estimateTx.from = tx.from;
        if (tx.to !== "" && tx.to !== "0x")
          estimateTx.to = tx.to;
        estimateTx.value = tx.value;
        estimateTx.nonce = tx.nonce;
        estimateTx.data = tx.data;
        return send("eth_estimateGas")(estimateTx)
          .then(usedGas => {
            return Map.merge(tx)({gas: Nat.div(Nat.mul(usedGas,"0x6"),"0x5")});
          });
      } else {
        return Promise.resolve(tx);
      }
    });

  const sendTransactionWithDefaults = tx =>
    addTransactionDefaults(tx)
      .then(tx => sendTransaction(removeEmptyTo(tx)));

  const callWithDefaults = (tx, block) =>
    addTransactionDefaults(tx)
      .then(tx => call(removeEmptyTo(tx), block || "latest"));

  const callMethodData = method => (...params) => {
    const methodSig = method.name + "(" + method.inputs.map(i => i.type).join(",") + ")";
    const methodHash = keccak256s(methodSig).slice(0,10);
    let encodedParams = params.map((param,i) => encodeABI(method.inputs[i].type, param));
    var headBlock = "0x";
    let dataBlock = "0x";
    for (var i = 0; i < encodedParams.length; ++i) {
      if (encodedParams[i].dynamic) {
        var dataLoc = encodedParams.length * 32 + Bytes.length(dataBlock);
        headBlock = Bytes.concat(headBlock, Bytes.pad(32, Nat.fromNumber(dataLoc)));
        dataBlock = Bytes.concat(dataBlock, encodedParams[i].data);
      } else {
        headBlock = Bytes.concat(headBlock, encodedParams[i].data);
      }
    }
    return Bytes.flatten([methodHash, headBlock, dataBlock]);
  }

  // Address, Address, ContractInterface -> Contract
  const contract = (from, address, contractInterface) => {
    let contract = {};
    contract._address = address;
    contract._from = from;
    contract.broadcast = {};
    contractInterface.forEach(method => {
      if (method && method.name) {
        const call = (type, value) => (...params) => {
          const transaction = {
            from: from,
            to: address,
            value: value,
            data: callMethodData(method)(...params)
          };
          return type === "data"
            ? Promise.resolve(transaction)
            : method.constant
              ? callWithDefaults(transaction)
              : sendTransactionWithDefaults(transaction)
                .then(type === "receipt" ? waitTransactionReceipt : (x => x));
        };
        contract[method.name] = call("receipt", "0x0");
        if (!method.constant) {
          contract[method.name+"_data"] = call("data", "0x0");
          contract[method.name+"_pay"] = value => (...params) => call("receipt", value)(...params);
          contract[method.name+"_pay_txHash"] = value => (...params) => call("txHash", value)(...params);
          contract[method.name+"_txHash"] = call("txHash", "0x0");
        }
      }
    });
    return contract;
  }

  // Address, Bytecode -> TxHash
  const deployBytecode_txHash = (from, code) =>
    sendTransactionWithDefaults({from: from, data: code, to: ""});

  // Address, Bytecode -> Receipt
  const deployBytecode = (from, code) =>
    deployBytecode_txHash(from,code)
      .then(waitTransactionReceipt);

  // Address, Bytecode, ContractInterface
  const deployBytecodeContract = (from, code, contractInterface) =>
    deployBytecode(from, code)
      .then(receipt => contract(from, receipt.contractAddress, contractInterface));
      
  // Address, String, Address -> Contract
  const solidityContract = (from, source, at) =>
    compileSolidity(source)
      .then(({info:{abiDefinition}}) => contract(from, at, abiDefinition));

  // Address, String -> TxHash
  const deploySolidity_txHash = (from, source) =>
    compileSolidity(source)
      .then(({code}) => deployBytecode_txHash(from, code));

  // Address, String -> Receipt
  const deploySolidity = (from, source) =>
    deploySolidity_txHash(from, source)
      .then(waitTransactionReceipt);

  // Address, String -> Contract
  const deploySolidityContract = (from, source) =>
    compileSolidity(source)
      .then(({code, info:{abiDefinition}}) =>
        deployBytecodeContract(from, code, abiDefinition));

  return {
    send,

    sendTransaction,
    sendRawTransaction,
    getTransactionReceipt,
    call,
    getBalance,
    accounts,

    waitTransactionReceipt,
    addTransactionDefaults,
    sendTransactionWithDefaults,
    callWithDefaults,
    callMethodData,

    contract,
    deployBytecode_txHash,
    deployBytecode,
    deployBytecodeContract,

    compileSolidity,
    solidityContract,
    deploySolidity_txHash,
    deploySolidity,
    deploySolidityContract
  }
}

module.exports = Api;
