This directory contains the workflows that would be used to handle SWIFT 
messages in the scenario where we are a middle man to DvP

SWIFT messaging goes as follows where S is the seller's bank, B is the buyer's bank, and D is the DON.
The SWIFT workflow DON is represented by D.

S -> B + D : sese.023 - Securities Settlement Transaction Instruction
B -> S + D: sese.024 - Securities Settlement Status Advices
// I think it can also be pacs.008, but the DON does not need to get this message anyways 
B -> S: pacs.009 - Financial Institution Credit Transfer
S -> B + D: pacs.002 - FI To FI Payment Status Report
D -> S + B: sese.025 - Securities Settlement Transaction Confirmation
