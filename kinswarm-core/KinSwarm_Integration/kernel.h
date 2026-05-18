#ifndef KERNEL_H
#define KERNEL_H
#include <stdint.h>
#include <stdbool.h>
typedef struct {
    uint8_t root_output[32];
    bool success;
} SettlementOutcome;
SettlementOutcome execute_settlement_anchor(const uint8_t* root_input, uint64_t amount, uint32_t worker_count);
#endif
