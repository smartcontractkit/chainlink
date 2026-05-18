OFFICIAL DISCLOSURE: THE KEEPER / THE ARCHITECT
INTENT: This intelligence is provided to Chainlink in good faith to facilitate robust infrastructure hardening. While these findings were previously submitted to and disregarded by the original protocols, they are handed over here under the assumption of professional integrity and collaborative advancement. This research is not public and remains anchored under the IUSC-1.0 Covenant to prevent the abuse or unauthorized replication encountered in prior disclosure attempts.

---

Cooldown Grief Attack in FXRP Redemption Enables Vault Liquidity Freeze via Agent Cycling
Submitted 11 months ago by @Pray4Love1 (Whitehat) for Flare Network
Report ID: 47098
Target: https://github.com/flare-foundation/bug-bounty/blob/main/smart-contracts.md
Impact(s): 
- Temporary freezing of funds
- Griefing (e.g. no profit motive for an attacker, but damage to the users or the protocol)
- Persistent state of queue exhaustion

Description:
The FXRP redemption mechanism in FAssets v1.1 enforces cooldown periods after redemptions to prevent abuse. However, it fails to prevent abuse via agent cycling — attackers can repeatedly register new disposable agents and submit minimal redemptions to exhaust all available agents, forcing the system into a blocked state. This causes denial of service for legitimate users, without breaching any explicit protocol rules.

Vulnerability Details:
The redemption system enforces a cooldown window per agent. Redemption queues rely on active agent availability. Cooldowns are time-based (block spans), but not identity-weighted.

An attacker can:
1. Register 20+ low-cost temporary agents (throwaway signers or bots).
2. Redeem a tiny FXRP amount from each, triggering cooldown.
3. Saturate the entire queue, making it impossible for legitimate agents to redeem.

There is no mechanism to prevent this since cooldown is agent-local and lacks global anti-cycling logic. Collateral slashing is not triggered by griefing behavior.

Impact Details:
This vulnerability enables a denial-of-service condition against the FXRP redemption queue by spamming redemptions through disposable agents. Legitimate redemptions become inaccessible, causing redemptions to fail or stall.
- Temporary freezing of funds for at least 24 hours.
- Vault liquidity becomes inaccessible for legitimate users until cooldowns expire.
- Risk scales with vault congestion and token volatility.

Proof of Concept:
```js
describe("Cooldown Grief Attack", function () { 
  it("uses disposable agents to saturate redemption queue", async function () { 
    const AgentFactory = await ethers.getContractFactory("MockCooldownAgent"); 
    const cooldown = 5; 
    const deployedAgents = [];

    // Register many temporary agents
    for (let i = 0; i < 20; i++) {  
      const agent = await AgentFactory.deploy();  
      await agent.registerAgent(cooldown);  
      deployedAgents.push(agent);
    }
    
    // Agents now enter cooldown
    for (let agent of deployedAgents) {  
      const canRedeem = await agent.redeem();  
      expect(canRedeem).to.be.false;
    }
    
    // Now try to redeem with a legit agent
    const legit = await AgentFactory.deploy();
    await legit.registerAgent(cooldown);
    const legitRedeem = await legit.redeem();
    expect(legitRedeem).to.be.false; // Cooldown grief successful
  }); 
});
