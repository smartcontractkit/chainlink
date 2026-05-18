#ifndef KERNEL_H
#define KERNEL_H

#include <stdint.h>
#include <stdbool.h>

typedef struct {
    bool success;
    uint8_t root_output[32];
} KernelOutcome;

KernelOutcome execute_settlement_batch(const uint8_t* merkle_root, uint64_t amount, uint32_t worker_count);

#endif
