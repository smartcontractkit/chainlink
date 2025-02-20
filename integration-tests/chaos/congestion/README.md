## Glossary
Congestion - block's property calculated by dividing the gas used in a block by its limit. Since EIP-1559 congestion
greater than 0.5 increases the base fee (minimal fee paid by a sender to consider a transaction valid).
Higher congestion correlates with transaction acceptance latency and higher tips (fees paid to incentivize block producers to include the transaction).
## Simulator
`congestion.Simulator` is a component that simulates spikes in chain activity by doing the following:
1. Overriding base fees to act according to configuration:
    1. Increases during ramp-up phase by `FeesIncreasePercent` from initial values;
    2. Maintains at an increased level during the plateau phase;
    3. Reduces back to baseline during cool down.
2. Producing transactions with increased tips according to current phases similar to base fee behavior. The number of produced transactions is adjusted to achieve target congestion.
