# Compiling

Using nightly:

    rustup override set nightly-2018-09-30
    rustup target add wasm32-unknown-unknown --toolchain nightly-2018-09-30
    cargo build --release

Shrinking the binary:

    cargo install --git https://github.com/alexcrichton/wasm-gc
    wasm-gc target/wasm32-unknown-unknown/release/multiply_wasm.wasm
