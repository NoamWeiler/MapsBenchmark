package internal

import (
	"MapsBenchmark/internal/internal_mutex_map"
	"MapsBenchmark/internal/mutex_map"
	"MapsBenchmark/internal/rw_map"
	"bufio"
	"fmt"
	randa "github.com/FTChinese/go-rest/rand"
	"log"
	"os"
	"sync"
	"time"
)

var (
	mutexMap  *mutex_map.MutexMap
	rwMap     *rw_map.RWMap
	syncMap   *sync.Map
	structMap *internal_mutex_map.InternalRWMutexMap
	mapSize   = 10000
	keysSlice []string
)

type maps interface {
	Add(k, v string)
	Get(k string) string
	Delete(k string)
}

func randString() string {
	return randa.String(20)
}

func initMap(m maps) {
	for i := 0; i < mapSize; i++ {
		m.Add(randString(), randString())
	}
}

func initSyncMap(sm *sync.Map) {
	for i := 0; i < mapSize; i++ {
		sm.Store(randString(), randString())
	}
}

func InitMaps() {
	rwMap = rw_map.NewMap()
	mutexMap = mutex_map.NewMap()
	structMap = internal_mutex_map.NewMap()
	syncMap = &sync.Map{}
}

func addTest(m maps) time.Duration {
	start := time.Now()
	initMap(m)
	return time.Since(start)
}

func addTestSyncMap(sm *sync.Map) time.Duration {
	start := time.Now()
	initSyncMap(sm)
	return time.Since(start)
}

// MapsAddTest run 10 times on each data structure to get better AVG
func MapsAddTest() {
	var sum1, sum2, sum3, sum4 time.Duration
	for i := 0; i < 10; i++ {
		InitMaps()

		sum1 += addTest(rwMap)
		sum2 += addTest(mutexMap)
		sum3 += addTest(structMap)
		sum4 += addTestSyncMap(syncMap)

	}
	fmt.Println("Add test results:")
	fmt.Println("rwMap:\t\t", sum1/10)
	fmt.Println("mutexMap:\t", sum2/10)
	fmt.Println("structMap:\t", sum3/10)
	fmt.Println("syncMap:\t", sum4/10)
}

func getPWD() string {
	var path, err = os.Getwd()
	if err != nil {
		log.Println(err)
	}
	return path
}

func initKeysSlice() []string {
	var KeysSlice []string
	f, err := os.Open(fmt.Sprintf("%s/%s", getPWD(), "internal/db/strings_db.txt"))
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s := scanner.Text()
		KeysSlice = append(KeysSlice, s)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return KeysSlice
}

// InitMapsFromFile add k,k to map from the file
func initMapsFromFile(sm *sync.Map, m maps) {
	for _, s := range keysSlice {
		if m != nil {
			m.Add(s, s)
		} else {
			sm.Store(s, s)
		}
	}
}

func getTest(sm *sync.Map, m maps) time.Duration {
	wg := &sync.WaitGroup{}

	//run the gouroutine on all strings from map
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, str := range keysSlice {
				if m == nil {
					if _, b := sm.Load(str); !b {
						fmt.Println(str)
					}
				} else {
					if m.Get(str) == "" {
						fmt.Println("Error, got nil for string:", str)
					}
				}
			}
		}()
	}

	now := time.Now()
	wg.Wait()
	return time.Since(now)
}

func getTestNoWG(sm *sync.Map, m maps) time.Duration {
	now := time.Now()

	for _, str := range keysSlice {
		if m == nil {
			_, _ = sm.Load(str)
		} else {
			_ = m.Get(str)
		}
	}

	return time.Since(now)
}

func getAndAddTest(sm *sync.Map, m maps) time.Duration {
	addChanRW := make(chan time.Duration, 1)
	addChansyncMap := make(chan time.Duration, 1)
	getChanRW := make(chan time.Duration, 1)
	getChansyncMap := make(chan time.Duration, 1)
	defer close(addChanRW)
	defer close(addChansyncMap)
	defer close(getChansyncMap)
	defer close(getChanRW)

	var readRW, readSync, addRW, addSync time.Duration
	if m == nil {
		go func(c chan<- time.Duration) {
			c <- addTestSyncMap(sm)
		}(addChansyncMap)

		go func(c chan<- time.Duration) {
			c <- getTestNoWG(sm, nil)
		}(getChansyncMap)
	} else {
		//add goroutine
		go func(c chan<- time.Duration) {
			c <- addTest(rwMap)
		}(addChanRW)

		go func(c chan<- time.Duration) {
			c <- getTestNoWG(nil, rwMap)
		}(getChanRW)
	}

	if m == nil {
		readSync = <-getChansyncMap
		addSync = <-addChansyncMap
		return readSync + addSync
	} else {
		readRW = <-getChanRW
		addRW = <-addChanRW
		return readRW + addRW
	}

}

func initialGetTests() {
	keysSlice = initKeysSlice()
	InitMaps()
	initMapsFromFile(nil, rwMap)
	initMapsFromFile(nil, mutexMap)
	initMapsFromFile(nil, structMap)
	initMapsFromFile(syncMap, nil)
}

func getTestStruct(withWG bool) {
	initialGetTests()
	var sum1, sum2, sum3, sum4 time.Duration
	var heading string
	if withWG {
		heading = "Get test with goroutines results:"
	} else {
		heading = "Get test without goroutines results:"
	}
	for i := 0; i < 10; i++ {
		if withWG {
			sum1 += getTest(nil, rwMap)
			sum2 += getTest(nil, mutexMap)
			sum3 += getTest(nil, structMap)
			sum4 += getTest(syncMap, nil)

		} else {
			sum1 += getTestNoWG(nil, rwMap)
			sum2 += getTestNoWG(nil, mutexMap)
			sum3 += getTestNoWG(nil, structMap)
			sum4 += getTestNoWG(syncMap, nil)
		}
	}

	fmt.Println(heading)
	fmt.Println("rwMap:\t\t", sum1/10)
	fmt.Println("mutexMap:\t", sum2/10)
	fmt.Println("structMap:\t", sum3/10)
	fmt.Println("syncMap:\t", sum4/10)
}

func getTestWithWG() {
	getTestStruct(true)
}

func getTestWithoutWG() {
	getTestStruct(false)
}

func MapsGetTest() {
	getTestWithWG()
	getTestWithoutWG()
}

func MapsGetAndAddTest() {
	initialGetTests()
	var sum1, sum2, sum3, sum4 time.Duration
	for i := 0; i < 10; i++ {
		sum1 += getAndAddTest(nil, rwMap)
		sum2 += getAndAddTest(nil, mutexMap)
		sum3 += getAndAddTest(nil, structMap)
		sum4 += getAndAddTest(syncMap, nil)
	}

	fmt.Println("Add test results:")
	fmt.Println("rwMap:\t\t", sum1/10)
	fmt.Println("mutexMap:\t", sum2/10)
	fmt.Println("structMap:\t", sum3/10)
	fmt.Println("syncMap:\t", sum4/10)

}
