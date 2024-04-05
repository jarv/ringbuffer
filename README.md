## Ring buffers in Go using channels

This repository has two implementations of a ring buffer in Go using channels.

1. Two channel implementation based on  https://tanzu.vmware.com/content/blog/a-channel-based-ring-buffer-in-go
2. Single channel implementation that uses a locking mutex

- The implementations in`bad/` contains two race conditions.
- The implementations in `googd/` have them resolved.
