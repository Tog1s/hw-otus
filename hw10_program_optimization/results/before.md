```bash
=== RUN   TestGetDomainStat_Time_And_Memory
    stats_optimization_test.go:46: time used: 669.810474ms / 300ms
    stats_optimization_test.go:47: memory used: 308Mb / 30Mb
    assertion_compare.go:332:
        	Error Trace:	stats_optimization_test.go:49
        	Error:      	"669810474" is not less than "300000000"
        	Test:       	TestGetDomainStat_Time_And_Memory
        	Messages:   	[the program is too slow]
--- FAIL: TestGetDomainStat_Time_And_Memory (29.88s)
FAIL
FAIL	github.com/tog1s/hw-otus/hw10_program_optimization	29.883s
FAIL
```