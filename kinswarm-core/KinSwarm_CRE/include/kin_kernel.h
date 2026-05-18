#ifndef KIN_KERNEL_H
#define KIN_KERNEL_H

#include <stdint.h>
#include <stdbool.h>

typedef struct {
    bool success;
    uint8_t message_hash[32];
} SettlementResult;

SettlementResult process_ledger_anchor(const uint8_t* merkle_root, uint64_t amount, uint32 worker_count);

#endif
