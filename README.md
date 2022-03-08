# MapsBenchmark

This project's aims is to demonstrate the difference between different map uses.</br>
Tested:</br>
* sync.Map
* map with mutex
* map with RWmutex
* map with rw mutex (as part of the value struct)

##Add Test
The data structres added pairs of random strings (used the sam string generator for all of them).</br>
Time measured from the beginning since the end f adding and were tested - time (for accurate avg)</br>
####The results for 100k pairs:</br>
* rwMap:           50.93645ms
* mutexMap:        49.047212ms
* structMap:       53.836724ms
* syncMap:         67.858754ms

####The results for 1M pairs:</br>
* rwMap:           359.360133ms
* mutexMap:        356.801833ms
* structMap:       358.063575ms
* syncMap:         608.426704ms




##Get Test
The data structures added pairs of strings from a generated big file (./internal/db)</br>
Time measured from the point were the waitingGroup is waiting for the goroutines to finish get command from the map.</br>
The result here were pretty suprising at first glance:

####Add test results:
* rwMap:           155.324908ms
* mutexMap:        140.764179ms
* structMap:       153.844012ms
* syncMap:         7.156266ms

The reason for these huge difference is that the sync.map there is a use of cache.


The second test wa without goroutine, in order o compare only the data structures' performance</br>
####The results:
* rwMap:           408.1µs
* mutexMap:        352.004µs
* structMap:       394.429µs
* syncMap:         524.512µs




