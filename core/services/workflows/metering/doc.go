/*
Metering / Billing manages a user balance across a workflow execution where multiple capabilities may wish to
deduct from the user balance. Additionally, a report will be generated and provided to both a billing service
and beholder.

Dependencies:
- Billing Service (provides: user balance, rate card, deduct/settle usage)
- Beholder (provides: report tracking)
- Capability spend ratios in capability config (where applicable)

Metering and Billing are considered separate topics where Metering is a value as either a limit or usage relating
to a capability. Billing relates to a user balance. Metered values can be converted to a balance value to be
deducted. To accomplish this conversion, a rate card is required and provided by a billing service.

Metering Mode:

Metering mode is a failure state where metering is still active and tracking capability spends, but a balance is
not being deducted and no balance updates will be sent to the billing service. Metering mode can be triggered by
multiple scenarios (refer to metering_test.go for more detail):

 1. Reserve switches to metering mode when
    a. billing client is nil or returns an error
    b. rate card from billing client contains invalid value (not parsable as decimal)
 2. ConvertToBalance switches to metering mode when a capability resource is not found in the billing rate card
 3. CreditToSpendingLimits switches to metering mode when
    a. a capability resource type does not exist in rate card
    b. a capability spends multiple resources without matching ratios

Metering mode can apply only once per metering report. Once in metering mode for a report, it will stay in metering
mode until completion. A new report will be started with metering mode off.

Once metering mode is triggered, an error will be written to the log indicating that metering mode is on and the reason
for metering mode. Once in metering mode, no more communication attempts will be made to the billing service for the
duration of the report. Beholder will be notified of a report in all cases.
*/
package metering
