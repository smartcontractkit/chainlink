This directory contains the workflows that would be used to handle SWIFT 
messages as-if the chain is acting as an escrow agent.

The DON will respond back on the chain's behalf with the SWIFT messages.

One advantage of this approach, assuming that escrows participate in this manner already, 
is that nothing is new for SWIFT during the transactions themselves.

The overall workflow for DvP, in terms of SWIFT messages is seen as follows,
where B is the buyer's band, S is the seller's bank, and D is the DON.

S -> B: sese.023 - Securities Settlement Transaction Instruction
B -> S: sese.024 - Securities Settlement Status Advices 
S -> D: 