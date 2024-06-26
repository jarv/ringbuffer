## Ring buffers in Go using channels

This repository has two implementations of a ring buffer in Go using channels.

See the accompanying blog post [Creating a ring buffer using Go Channels for Server-Sent Events](https://jarv.org/posts/go-channels-sse/).

1. Two channel implementation based on https://tanzu.vmware.com/content/blog/a-channel-based-ring-buffer-in-go
2. Single channel implementation that uses a locking mutex

- `bad/` contains the implementations with race conditions.
- `good/` has the race conditions resolved.


### Unit tests

```
go test ./...
```

### Benchmarks

```
cd good
go test -bench=. ./...
```
