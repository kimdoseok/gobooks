#!/bin/sh
cargo install trunk
rustup target add wasm32-unknown-unknown
cargo new yewapp
cd yewapp
cargo run


