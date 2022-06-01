import { ethers } from "hardhat"
import { GovernanceAutomator } from "../../typechain/GovernanceAutomator"
import { expect } from "chai"
import { deployMockContract, MockContract } from '@ethereum-waffle/mock-contract'
import IGovernanceABI from "../../abi/src/v0.8/dev/interfaces/IGovernance.sol/IGovernance.json"
import IGovernanceTokenABI from "../../abi/src/v0.8/dev/interfaces/IGovernanceToken.sol/IGovernanceToken.json"


let governanceMock: MockContract
let governanceMock_token: MockContract
let governanceAutomator: GovernanceAutomator
let pAddress = "0xc26d7EF337e01a5cC5498D3cc2ff0610761ae637"

enum ProposalState {Pending, Canceled, Active, Failed, Succeeded, Queued, Expired, Executed}
enum Action {QUEUE, EXECUTE, CANCEL, INDEX}

async function setMockedProposals(proposalStates: number[]){
    await governanceMock.mock.proposalCount.returns(proposalStates.length)
    for (let idx = 0; idx < proposalStates.length; idx++) {
        const proposalState = proposalStates[idx];
        await governanceMock.mock.state.withArgs(idx + 1).returns(proposalState)
    }
}

let test_proposal = {
    id: 1,
    proposer: pAddress,
    eta: 123,
    targets: [],
    values: [],
    signatures: [],
    calldatas: [],
    startBlock: 1,
    endBlock: 1,
    forVotes: 1,
    againstVotes: 1,
    canceled: false,
    executed: false
}

async function getCancelProposal(){
    await governanceMock.mock.proposals.withArgs(1).returns(test_proposal)
    await governanceMock.mock.proposalThreshold.withArgs().returns(2)
    await governanceMock_token.mock.getPriorVotes.returns(1)
}

async function getValidProposal(){
    await governanceMock.mock.proposals.withArgs(1).returns(test_proposal)
    await governanceMock.mock.proposalThreshold.withArgs().returns(1)
    await governanceMock_token.mock.getPriorVotes.returns(2)

}

describe.only("GovernanceAutomator Contract", ()=>{

    beforeEach(async ()=>{
        const wallet = (await ethers.getSigners())[0]
        const contractAbi = IGovernanceABI
        const tokenContractAbi = IGovernanceTokenABI
        governanceMock = await deployMockContract(wallet as any, contractAbi);
        governanceMock_token = await deployMockContract(wallet as any, tokenContractAbi);
        var autoFactory = await ethers.getContractFactory("GovernanceAutomator")
        governanceAutomator = await autoFactory.deploy(governanceMock.address, 1, governanceMock_token.address)
    })
    describe("checkUpkeep", ()=>{

        it("Returns true for eligible state and false for ineligible state", async ()=>{
            await getValidProposal()

            await setMockedProposals([ProposalState.Pending])
            let result = await governanceAutomator.checkUpkeep("0x")
            expect(result.upkeepNeeded).to.be.false

            await setMockedProposals([ProposalState.Canceled])
            result = await governanceAutomator.checkUpkeep("0x")
            expect(result.upkeepNeeded).to.be.false

            await setMockedProposals([ProposalState.Active])
            result = await governanceAutomator.checkUpkeep("0x")
            expect(result.upkeepNeeded).to.be.false

            await setMockedProposals([ProposalState.Failed])
            result = await governanceAutomator.checkUpkeep("0x")
            expect(result.upkeepNeeded).to.be.false

            await setMockedProposals([ProposalState.Succeeded])
            result = await governanceAutomator.checkUpkeep("0x")
            expect(result.upkeepNeeded).to.be.true
            expect(result.performData).to.eq(ethers.utils.defaultAbiCoder.encode(["uint8","uint256"],[Action.QUEUE, 1]))

            await setMockedProposals([ProposalState.Queued])
            result = await governanceAutomator.checkUpkeep("0x")
            expect(result.upkeepNeeded).to.be.true
            expect(result.performData).to.eq(ethers.utils.defaultAbiCoder.encode(["uint8","uint256"],[Action.EXECUTE, 1]))

            await setMockedProposals([ProposalState.Expired])
            result = await governanceAutomator.checkUpkeep("0x")
            expect(result.upkeepNeeded).to.be.false

            await setMockedProposals([ProposalState.Executed])
            result = await governanceAutomator.checkUpkeep("0x")
            expect(result.upkeepNeeded).to.be.false
        })

        it("Returns true if index requires an update", async ()=>{
            await getValidProposal()
            await setMockedProposals([ProposalState.Canceled,ProposalState.Pending,ProposalState.Pending])
            let result = await governanceAutomator.checkUpkeep("0x")
            expect(result.upkeepNeeded).to.be.true
            expect(result.performData).to.eq(ethers.utils.defaultAbiCoder.encode(["uint8","uint256"],[Action.INDEX, 2]))
        })

        it("Returns true if index requires an update and index is last item in array", async ()=>{
            await getValidProposal()
            await setMockedProposals([ProposalState.Canceled,ProposalState.Canceled,ProposalState.Canceled])
            let result = await governanceAutomator.checkUpkeep("0x")
            expect(result.upkeepNeeded).to.be.true
            expect(result.performData).to.eq(ethers.utils.defaultAbiCoder.encode(["uint8","uint256"],[Action.INDEX, 3]))
        })

        it("Returns true if index requires an update and index is last item in array", async ()=>{
            await getValidProposal()
            await setMockedProposals([ProposalState.Canceled,ProposalState.Pending,ProposalState.Canceled, ProposalState.Pending])
            let result = await governanceAutomator.checkUpkeep("0x")
            expect(result.upkeepNeeded).to.be.true
            expect(result.performData).to.eq(ethers.utils.defaultAbiCoder.encode(["uint8","uint256"],[Action.INDEX, 2]))
        })



    }) 

    describe("performUpkeep", ()=>{
        it("Call cancel() function on gov contract in performUpkeep w/ proposal in 'Pending' State' (should revert)", async ()=>{
            await setMockedProposals([ProposalState.Pending]); 
            await getCancelProposal()
            let result = await governanceAutomator.checkUpkeep("0x")
            await governanceMock.mock.cancel.withArgs(1).returns()
            await expect(governanceAutomator.performUpkeep(result.performData))
        })
        it("Call queue() function on gov contract in performUpkeep w/ proposal in 'Pending' State' (should revert)", async ()=>{
            await getValidProposal()
            await setMockedProposals([ProposalState.Pending]); 
            let result = await governanceAutomator.checkUpkeep("0x")
            await governanceMock.mock.queue.withArgs(1).returns()
            await expect(governanceAutomator.performUpkeep(result.performData)).to.be.reverted
        })
        it("Call queue() function on gov contract in performUpkeep w/ proposal in 'Canceled' State' (should revert)", async ()=>{
            await getValidProposal()
            await setMockedProposals([ProposalState.Canceled]);
            let result = await governanceAutomator.checkUpkeep("0x")
            await governanceMock.mock.queue.withArgs(1).returns()
            await expect(governanceAutomator.performUpkeep(result.performData)).to.be.reverted
        })  
        it("Call queue() function on gov contract in performUpkeep w/ proposal in 'Active' State' (should revert)", async ()=>{
            await getValidProposal()
            await setMockedProposals([ProposalState.Active]);
            let result = await governanceAutomator.checkUpkeep("0x")
            await governanceMock.mock.queue.withArgs(1).returns()
            await expect(governanceAutomator.performUpkeep(result.performData)).to.be.reverted
        }) 
        it("Call queue() function on gov contract in performUpkeep w/ proposal in 'Failed' State' (should revert)", async ()=>{
            await getValidProposal()
            await setMockedProposals([ProposalState.Failed]);
            let result = await governanceAutomator.checkUpkeep("0x")
            await governanceMock.mock.queue.withArgs(1).returns()
            await expect(governanceAutomator.performUpkeep(result.performData)).to.be.reverted
        }) 
        it("Call queue() function on gov contract in performUpkeep w/ proposal in 'Succeeded' State' (should succeed)", async ()=>{
            await getValidProposal()
            await setMockedProposals([ProposalState.Succeeded]);
            let result = await governanceAutomator.checkUpkeep("0x")
            await governanceMock.mock.queue.withArgs(1).returns()
            await governanceAutomator.performUpkeep(result.performData)
        }) 
        it("Call queue() function on gov contract in performUpkeep w/ proposal in 'Queued' State' (should revert)", async ()=>{
            await getValidProposal()
            await setMockedProposals([ProposalState.Queued]); 
            let result = await governanceAutomator.checkUpkeep("0x")
            await governanceMock.mock.queue.withArgs(1).returns()
            await expect(governanceAutomator.performUpkeep(result.performData)).to.be.reverted
        }) 
        it("Call queue() function on gov contract in performUpkeep w/ proposal in 'Expired' State' (should revert)", async ()=>{
            await getValidProposal()
            await setMockedProposals([ProposalState.Expired]); 
            let result = await governanceAutomator.checkUpkeep("0x")
            await governanceMock.mock.queue.withArgs(1).returns()
            await expect(governanceAutomator.performUpkeep(result.performData)).to.be.reverted
        }) 
        it("Call queue() function on gov contract in performUpkeep w/ proposal in 'Executed' State' (should revert)", async ()=>{
            await getValidProposal()
            await setMockedProposals([ProposalState.Executed]); 
            let result = await governanceAutomator.checkUpkeep("0x")
            await governanceMock.mock.queue.withArgs(1).returns()
            await expect(governanceAutomator.performUpkeep(result.performData)).to.be.reverted
        }) 

        // Executed checks

        it("Call execute() function on gov contract in performUpkeep w/ proposal in 'Pending' State' (should revert)", async ()=>{
            await getValidProposal()
            await setMockedProposals([ProposalState.Pending]); 
            let result = await governanceAutomator.checkUpkeep("0x")
            await governanceMock.mock.execute.withArgs(1).returns()
            await expect(governanceAutomator.performUpkeep(result.performData)).to.be.reverted
        })
        it("Call execute() function on gov contract in performUpkeep w/ proposal in 'Canceled' State' (should revert)", async ()=>{
            await getValidProposal()
            await setMockedProposals([ProposalState.Canceled]); 
            let result = await governanceAutomator.checkUpkeep("0x")
            await governanceMock.mock.execute.withArgs(1).returns()
            await expect(governanceAutomator.performUpkeep(result.performData)).to.be.reverted
        })  
        it("Call execute() function on gov contract in performUpkeep w/ proposal in 'Active' State' (should revert)", async ()=>{
            await getValidProposal()
            await setMockedProposals([ProposalState.Active]);
            let result = await governanceAutomator.checkUpkeep("0x")
            await governanceMock.mock.execute.withArgs(1).returns()
            await expect(governanceAutomator.performUpkeep(result.performData)).to.be.reverted
        }) 
        it("Call execute() function on gov contract in performUpkeep w/ proposal in 'Failed' State' (should revert)", async ()=>{
            await getValidProposal()
            await setMockedProposals([ProposalState.Failed]);
            let result = await governanceAutomator.checkUpkeep("0x")
            await governanceMock.mock.execute.withArgs(1).returns()
            await expect(governanceAutomator.performUpkeep(result.performData)).to.be.reverted
        }) 
        it("Call execute() function on gov contract in performUpkeep w/ proposal in 'Succeeded' State' (should revert)", async ()=>{
            await getValidProposal()
            await setMockedProposals([ProposalState.Succeeded]);
            let result = await governanceAutomator.checkUpkeep("0x")
            await governanceMock.mock.execute.withArgs(1).returns()
            await expect(governanceAutomator.performUpkeep(result.performData)).to.be.reverted  
        }) 
        it("Call execute() function on gov contract in performUpkeep w/ proposal in 'Queued' State' (should succeed)", async ()=>{
            await getValidProposal()
            await setMockedProposals([ProposalState.Queued]); 
            let result = await governanceAutomator.checkUpkeep("0x")
            await governanceMock.mock.execute.withArgs(1).returns()
            await governanceAutomator.performUpkeep(result.performData)
        }) 
        it("Call execute() function on gov contract in performUpkeep w/ proposal in 'Expired' State' (should revert)", async ()=>{
            await getValidProposal()
            await setMockedProposals([ProposalState.Expired]); 
            let result = await governanceAutomator.checkUpkeep("0x")
            await governanceMock.mock.execute.withArgs(1).returns()
            await expect(governanceAutomator.performUpkeep(result.performData)).to.be.reverted
        }) 
        it("Call execute() function on gov contract in performUpkeep w/ proposal in 'Executed' State' (should revert)", async ()=>{
            await getValidProposal()
            await setMockedProposals([ProposalState.Executed]); 
            let result = await governanceAutomator.checkUpkeep("0x")
            await governanceMock.mock.execute.withArgs(1).returns()
            await expect(governanceAutomator.performUpkeep(result.performData)).to.be.reverted
        }) 
    })

})
