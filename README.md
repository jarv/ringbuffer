## Ring buffers in Go using channels

This repository has two implementations of a ring buffer in Go using channels.

1. Two channel implementation inspired by  https://tanzu.vmware.com/content/blog/a-channel-based-ring-buffer-in-go
2. Single channel implementation that uses a mutex to prevent race conditions

Included is the script `bench` to run benchmark tests against the two implementations.
