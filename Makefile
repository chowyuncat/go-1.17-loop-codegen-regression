GO116=/opt/homebrew/opt/go@1.16/bin/go
GO118=/opt/homebrew/opt/go@1.18/bin/go
BENCH_N=30000

# macOS only: uncomment to force execution on the efficiency cores of an M1
# XPROG=taskpolicy -c background

GOFLAGS = build -gcflags='-l'

all: loop_116 loop_118

clean:
	rm loop_116 loop_118; true

loop_116: loop.go
	$(GO116) $(GOFLAGS) -o $@ $<

loop_118: loop.go
	$(GO118) $(GOFLAGS) -o $@ $<

bench: loop_116 loop_118
	time $(XPROG) ./loop_116 $(BENCH_N) pre
	time $(XPROG) ./loop_116 $(BENCH_N) post
	echo
	time $(XPROG) ./loop_118 $(BENCH_N) pre
	time $(XPROG) ./loop_118 $(BENCH_N) post
