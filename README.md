# tetris

A trivial implementation of Tetris in Golang

## Debug

### Profiling

To profile CPU usage, run the following command.

```bash
go run main.go -cpuprofile cpu.prof
```

or if you are interested in memory usage,

```bash
go run main.go -memprofile mem.prof
```

After obtaining the profile, analyze it with the following command.

```bash
$ go tool pprof cpu.prof  # or mem.prof
File: main
Type: cpu
Time: May 10, 2024 at 7:19pm (JST)
Duration: 4.13s, Total samples = 3.18s (77.04%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top 10
Showing nodes accounting for 3.10s, 97.48% of 3.18s total
Dropped 59 nodes (cum <= 0.02s)
Showing top 10 nodes out of 93
      flat  flat%   sum%        cum   cum%
     0.86s 27.04% 27.04%      0.86s 27.04%  <unknown>
     0.78s 24.53% 51.57%      0.78s 24.53%  runtime.pthread_cond_wait
     0.77s 24.21% 75.79%      0.77s 24.21%  runtime.cgocall
     0.23s  7.23% 83.02%      0.23s  7.23%  runtime.pthread_cond_signal
     0.23s  7.23% 90.25%      0.23s  7.23%  runtime.pthread_kill
     0.13s  4.09% 94.34%      0.13s  4.09%  runtime.usleep
     0.05s  1.57% 95.91%      0.05s  1.57%  runtime.pthread_cond_timedwait_relative_np
     0.02s  0.63% 96.54%      0.02s  0.63%  runtime.kevent
     0.02s  0.63% 97.17%      0.04s  1.26%  runtime.mallocgc
     0.01s  0.31% 97.48%      0.03s  0.94%  runtime.gcBgMarkWorker
```

You can visualize the result of profiling.

```bash
go tool pprof -png mem.prof > out.png  # or cpu.prof
```

## Deployment

```zsh
env GOOS=js GOARCH=wasm go build -o tetris.wasm github.com/okayama-daiki/tetris
```
