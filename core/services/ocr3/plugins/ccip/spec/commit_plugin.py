#
# High-level Python specification for the CCIP OCR3 Commit Plugin.
#
# This specification aims to provide a clear and comprehensive understanding
# of the plugin's functionality. It is highly recommended for engineers working on CCIP
# to familiarize themselves with this specification prior to reading the
# corresponding Go implementation.
#
# NOTE: Even though the specification is written in a high-level programming language, it's purpose
# is not to be executed. It is meant to be just a reference for the Go implementation.
#
from dataclasses import dataclass
from typing import List, Dict

ChainSelector = int

@dataclass
class Interval:
    min: int
    max: int

@dataclass
class Message:
    seq_nr: int
    message_id: bytes    # a unique message identifier computed on the source chain
    message_hash: bytes  # hash of message body computed on the destination chain and used on merkle tree
    # TODO:

@dataclass
class Commit:
    interval: Interval
    root: bytes

@dataclass
class CommitOutcome:
    latest_committed_seq_nums: Dict[ChainSelector, int]
    commits: Dict[ChainSelector, Commit]
    token_prices: Dict[str, int]
    gas_prices: Dict[ChainSelector, int]

@dataclass
class CommitObservation:
    latest_committed_seq_nums: Dict[ChainSelector, int]
    new_msgs: Dict[ChainSelector, List[Message]]
    token_prices: Dict[str, int]
    gas_prices: Dict[ChainSelector, int]
    f_chain: Dict[ChainSelector, int]

@dataclass
class CommitConfig:
    oracle: int # our own observer
    dest_chain: ChainSelector
    f_chain: Dict[ChainSelector, int]
    # oracleIndex -> supported chains
    oracle_info: Dict[int, Dict[ChainSelector, bool]]
    priced_tokens: List[str]

class CommitPlugin:
    def __init__(self):
         self.cfg = CommitConfig(
            oracle=1,
            dest_chain=10,
            f_chain={1: 2, 2: 3, 10: 3},
            oracle_info={
                0: {1: True, 2: True, 10: True},
                # TODO: other oracles
            },
            # TODO: will likely need aggregator address as well to
            # actually get the price.
            priced_tokens=["tokenA", "tokenB"],
         )
         self.keep_cfg_in_sync()

    def get_token_prices(self):
        # Read token prices which are required for the destination chain.
        # We only read them if we have the capability to read from the price chain (e.g. arbitrum)
        pass

    def get_gas_prices(self):
        # Read all gas prices for the chains we support.
        pass

    def query(self):
        pass

    def observation(self, previous_outcome: CommitOutcome) -> CommitObservation:
        # max_committed_seq_nr={sourceChainA: 10, sourceChainB: 20,...}
        # Provided by the nodes that can read from the destination on the previous round.
        # Observe msgs for our supported chains since the prev outcome.
        new_msgs = {}
        for (chain, seq_num) in previous_outcome.latest_committed_seq_nums:
            if chain in self.cfg.oracle_info[self.cfg.oracle]:
                msgs = self.onRamp(chain).get_msgs(chain, start=seq_num+1, limit=256)
                for msg in msgs:
                    msg.message_hash = msg.compute_hash()
                new_msgs[chain] = msgs

        # Observe token prices. {token: price}
        token_prices = self.get_token_prices()

        # Observe gas prices. {chain: gasPrice}
        # TODO: Should be a way to combine the loops over support chains for gas prices and new messages.
        gas_prices = self.get_gas_prices()

        # Observe fChain for each chain. {chain: f_chain}
        # We observe this because configuration changes may be detected at different times by different nodes.
        # We always use the configuration which is seen by a majority of nodes.
        f_chain = self.cfg.f_chain

        # If we support the destination chain, then we contribute an observation of the max committed seq nums.
        # We use these in outcome to filter out messages that have already been committed.
        latest_committed_seq_nums = {}
        if self.cfg.dest_chain in self.cfg.oracle_info[self.cfg.oracle]:
            latest_committed_seq_nums = self.offRamp.latest_committed_seq_nums()

        return CommitObservation(latest_committed_seq_nums, new_msgs, token_prices, gas_prices, f_chain)


    def validate_observation(self, attributed_observation):
        observation = attributed_observation.observation
        oracle = attributed_observation.oracle

        # Only accept dest observations from nodes that support the dest chain
        if observation.latest_committed_seq_nums is not None:
            assert self.cfg.dest_chain in self.cfg.oracle_info[oracle]

        # Only accept source observations from nodes which support those sources.
        msg_ids = set()
        msg_hashes = set()
        for (chain, msgs) in observation.new_msgs.items():
            assert(chain in self.cfg.oracle_info[oracle])
            # Don't allow duplicates of (chain, seqNr), (id) and (hash). Required to prevent double counting.
            assert(len(msgs) == len(set([msg.seq_num for msg in msgs])))            
            for msg in msgs:
                assert msg.message_id not in msg_ids
                assert msg.message_hash not in msg_hashes
                msg_ids.add(msg.message_id)
                msg_hashes.add(msg.message_hash)

    def observation_quorum(self):
        return "2F+1"

    def outcome(self, observations: List[CommitObservation])->CommitOutcome:
        f_chain = consensus_f_chain(observations)
        latest_committed_seq_nums = consensus_latest_committed_seq_nums(observations, f_chain)

        # all_msgs contains all messages from all observations, grouped by source chain
        all_msgs = [observation.new_msgs for observation in observations].group_by_source_chain()

        commits = {} # { chain: (root, min_seq_num, max_seq_num) }
        for (chain, msgs) in all_msgs:
            # Keep only msgs with seq nums greater than the consensus max commited seq nums.
            # Note right after a report has been submitted, we'll expect those same messages
            # to appear in the next observation, because the message observations are built
            # on the previous max committed seq nums.
            msgs = [msg for msg in msgs if msg.seq_num > latest_committed_seq_nums[chain]]

            msgs_by_seq_num = msgs.group_by_seq_num() # { 423: [0x1, 0x1, 0x2] }
                                                      # 2 nodes say that msg hash is 0x1 and 1 node says it's 0x2
                                                      # if different hashes have the same number of votes, we select the
                                                      # hash with the lowest lexicographic order

            msg_hashes = { seq_num: elem_most_occurrences(hashes) for (seq_num, hashes) in msgs_by_seq_num.items() }
            for (seq_num, hash) in msg_hashes.items(): # require at least 2f+1 observations of the voted hash
                assert(msgs_by_seq_num[seq_num].count(hash) >= 2*f_chain[chain]+1)

            msgs_for_tree = [] # [ (seq_num, hash) ]
            for (seq_num, hash) in msg_hashes.ordered_by_seq_num():
                if len(msgs_for_tree) > 0 and msgs_for_tree[-1].seq_num+1 != seq_num:
                    break # gap in sequence numbers, stop here
                msgs_for_tree.append((seq_num, hash))

            commits[chain] = Commit(root=build_merkle_tree(msgs_for_tree), interval=Interval(min=msgs_for_tree[0].seq_num, max=msgs_for_tree[-1].seq_num))

        # TODO: we only want to put token/gas prices onchain
        # on a regular cadence unless huge deviation.
        token_prices = { tk: median(prices) for (tk, prices) in observations.group_token_prices_by_token() }
        gas_prices = { chain: median(prices) for (chain, prices) in observations.group_gas_prices_by_chain() }

        return CommitOutcome(latest_committed_seq_nums=latest_committed_seq_nums, commits=commits, token_price=token_prices, gas_prices=gas_prices)

    def reports(self, outcome):
        report = report_from_outcome(outcome)
        encoded = report.chain_encode() # abi_encode for evm chains
        return [encoded]

    def should_accept(self, report):
        if len(report) == 0 or self.validate_report(report):
            return False

    def should_transmit(self, report):
        if not self.is_writer():
            return False

        if len(report) == 0 or not self.validate_report(report):
            return False

        on_chain_seq_nums = self.offRamp.get_sequence_numbers()
        for (chain, tree) in report.trees():
            if not (on_chain_seq_nums[chain]+1 == tree.min_seq_num):
                return False

        return True

    def validate_report(self, report):
        pass

    def keep_cfg_in_sync(self):
        # Polling the configuration on the on-chain contract.
        # When the config is updated on-chain, updates the plugin's local copy to the most recent version.
        pass

def consensus_f_chain(observations):
    f_chain_votes = observations["f_chain"].group_by_chain() # { chainA: [1, 1, 16, 16, 16, 16] }
    return { ch: elem_most_occurrences(fs) for (ch, fs) in f_chain_votes.items() } # { chainA: 16 }

def consensus_latest_committed_seq_nums(observations, f_chains):
    all_latest_committed_seq_nums = {}
    for observation in observations:
        for (chain, seq_num) in observation.latest_committed_seq_nums.items():
            if chain not in all_latest_committed_seq_nums:
                all_latest_committed_seq_nums[chain] = []
            all_latest_committed_seq_nums[chain].append(seq_num)

    latest_committed_seq_nums_consensus = {}
    # { chainA: [4, 5, 5, 5, 5, 6, 6] }
    for (chain, latest_committed_seq_nums) in all_latest_committed_seq_nums.items():
        if len(latest_committed_seq_nums) >= 2*f_chains[chain]+1:
             # 2f+1 = 2*5+1 = 11
             latest_committed_seq_nums_consensus[chain] = sorted(latest_committed_seq_nums)[f_chains[chain]]# with f=4 { chainA: 5 }
    return latest_committed_seq_nums_consensus

def elem_most_occurrences(lst):
    pass

def build_merkle_tree(messages):
    pass

def median(lst):
    pass

def report_from_outcome(outcome: CommitOutcome)->bytes:
    pass
