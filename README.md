# go-boggle

Learning Golang via Solving a Boggle Board

## Building / Running

```bash
go build ./cmd/boggle
./boggle
```

## Testing

## Benchmarking

Benchmarking is a rough guess on how well certain parts of the application will perform.
Included in the repo are the baseline benchmarks as `baseline.out`.
Adjustments to the algorithm can be compared to these baselines.
If changes happen to these algorithms, the baseline should be updated.

```bash
cd internal/gameboard
go test -benchtime 3s -benchmem -run none -bench . > bench.out
benchcmp baseline.out bench.out
```

```bash
cd internal/lexicon
go test -benchtime 3s -benchmem -run none -bench . > bench.out
benchcmp baseline.out bench.out
```
