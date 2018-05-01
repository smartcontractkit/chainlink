function checkAllBalances() {
  var totalBal = 0;
  for (var acctNum in eth.accounts) {
    var acct = eth.accounts[acctNum];
    var acctBal = web3.fromWei(eth.getBalance(acct), "ether");
    totalBal += parseFloat(acctBal);
    console.log("  eth.accounts[" + acctNum + "]: \t" + acct + " \tbalance: " + acctBal + " ether");
  }
  console.log("  Total balance: " + totalBal + " ether");
};

function fundAccount() {
  return eth.sendTransaction({
    from:eth.accounts[0],
    to:eth.accounts[1],
    value: web3.toWei(1000, "ether")
  });
};
