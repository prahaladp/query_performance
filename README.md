# query_performance\

This tool measures the query performance for time series database.\
Assumptions : the data has already been stored in the database\

Usage\

Usage of ./query_performance:\
  -dbport int\
        db port to connect to (default 5432)\
  -file string
        filename containing the query parameters
  -host string
        destination hostname for db (default "localhost")
  -password string
        db user password (default "password")
  -tablename string
        table to query (default "cpu_usage")
  -username string
        db user name (default "postgres")
  -workers int
        number of worker threads (default 1)
      
The filename is the input file which contains the time slots per host
for which data is required.
 
The file is in the following format
 hostname,start_time,end_time
host_000008,2017-01-01 08:59:22,2017-01-01 09:59:22
host_000001,2017-01-02 13:02:02,2017-01-02 14:02:02
...

Sample Output 
$ ./query_performance  --file ./testdata/query_usage.csv --workers 1 --dbport 5432
---------------------------------
 number of samples = 200
 mean duration = 30.49438ms
 median duration = 29.855842ms
---------------------------------
Hostname : host_000007
        2017-01-02 17:43:00, 39.34, 67.60
        2017-01-02 17:44:00, 3.28, 57.21
        2017-01-02 17:45:00, 8.25, 83.98
...
Hostname : host_000004
        2017-01-01 08:58:00, 8.31, 92.96
        2017-01-01 08:59:00, 2.73, 64.17
        2017-01-01 09:00:00, 4.09, 90.52
...
Hostname : host_000001
        2017-01-02 13:02:00, 13.26, 98.39
        2017-01-02 13:03:00, 27.63, 92.11
        2017-01-02 13:04:00, 3.22, 96.99
....
        2017-01-02 15:35:00, 10.31, 68.40
        2017-01-02 15:36:00, 0.65, 84.09
        2017-01-02 15:37:00, 18.50, 82.74


Sample output from the tool
NOTE : this is obtained from multiple runs by changing the worker threads

Workers Mean              Median
---------------------------------------
1         31.626721ms       30.928894ms
2         31.92066ms        31.081195ms
4         56.676925ms       60.135168ms
8         113.060989ms    122.369957ms
16      148.308213ms      156.706489ms

The code should be documented to allow initial understanding of the data structures.
Most of the important data structures are in queryParams.go

Tests :
go test -v
=== RUN   TestMedian
--- PASS: TestMedian (0.00s)
=== RUN   TestBasicFile
--- PASS: TestBasicFile (0.00s)
=== RUN   TestSimpleDbOps
--- PASS: TestSimpleDbOps (0.01s)
=== RUN   TestComponentTest1
number of time samples = 1
--- PASS: TestComponentTest1 (0.02s)
=== RUN   TestBasic
--- PASS: TestBasic (0.00s)
=== RUN   TestQueryResults1
--- PASS: TestQueryResults1 (0.00s)
=== RUN   TestSimpleWorker
---------------------------------
 number of samples = 3
 mean duration = 862.739µs
 median duration = 621.199µs
---------------------------------
--- PASS: TestSimpleWorker (0.02s)
PASS
ok

NOTE : (obviously) needs more tests

$ go test -cover
PASS
coverage: 62.4% of statements
ok





        
