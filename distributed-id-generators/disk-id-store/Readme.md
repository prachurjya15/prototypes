# Disk-backed Distributed ID Generator

A monotonically increasing ID generator service backed by disk persistence, exposed via Go `net/rpc`.  
Multiple services request unique IDs over RPC. The generator uses a mutex for concurrency safety and periodically fsyncs to disk.

---

## Throughput Comparison

> **Services (fixed): 7** | Each service makes 10 ID generation calls  
> Latency measured from RPC call start to response received

| Strategy | No. of Services | Workers (goroutines) | Avg Time per ID (ms) |
|---|:---:|:---:|:---:|
| Write to disk on every call (no batching) | 7 | 7 | ~5.30 |
| Batched write (fsync every 10 IDs, Buffer=10) | 7 | 7 | ~0.34 |
| Batched write (fsync every 100 IDs, Buffer=100) | 7 | 7 | ~0.24 |
