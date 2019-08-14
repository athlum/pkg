package concurrentMap

import (
	"hash"
	"hash/adler32"
	"hash/fnv"
	"sync"
	"sync/atomic"
	"testing"
)

var (
	m *ConcurrentMap

	sm *sync.Map

	omp *atomic.Value

	om map[string]int

	l = []byte("afwopjapwfnapiwfjawpdokawdmapvbnoanfpaiwdkao,wdamo nvaonbpjfpawodmamcpoa,ponvpownfpanwfpgva owmfpoc,apmxaowmcpoanvpbianwifjapokwdpoamxpdoanmpocnapvnpianbipahbeinbpqpmpwdc,ao,mcpoam pboanpie biqemfpicmqjwpidjap Zpa epfimpeafmpoanvpamvpe nbpienfpiajpfockmavmpaonegvpiaifcmnaiopnscvia")

	s = []byte("bwm_lpd.eleyunying")
)

func init() {
	m = AdvanceNew(adler32.New, 16)
	m.Set("a", 1)
	m.Set("bwm_lpd.eleyunying", 1)

	sm = &sync.Map{}
	sm.Store("a", 1)
	sm.Store("bwm_lpd.eleyunying", 1)

	om = make(map[string]int)
	om["a"] = 1
	om["bwm_lpd.eleyunying"] = 1

	omp = &atomic.Value{}
	omp.Store(om)
}

func TestRead(t *testing.T) {
	wg := &sync.WaitGroup{}
	N := 10000000
	wg.Add(N)
	n := 20000
	if n > N {
		n = N
	}
	for i := 0; i < n; i++ {
		go func(wg *sync.WaitGroup) {
			for j := 0; j < N/n; j++ {
				m.Has("bwm_lpd.eleyunying")
				wg.Done()
			}
		}(wg)
	}
	wg.Wait()
}

func TestReadSync(t *testing.T) {
	wg := &sync.WaitGroup{}
	N := 10000000
	wg.Add(N)
	n := 20000
	if n > N {
		n = N
	}
	for i := 0; i < n; i++ {
		go func(wg *sync.WaitGroup) {
			for j := 0; j < N/n; j++ {
				sm.Load("bwm_lpd.eleyunying")
				wg.Done()
			}
		}(wg)
	}
	wg.Wait()
}

func TestReadOriginal(t *testing.T) {
	wg := &sync.WaitGroup{}
	N := 10000000
	wg.Add(N)
	n := 20000
	if n > N {
		n = N
	}
	for i := 0; i < n; i++ {
		go func(wg *sync.WaitGroup) {
			for j := 0; j < N/n; j++ {
				if _, ok := omp.Load().(map[string]int)["bwm_lpd.eleyunying"]; !ok {
					continue
				}
				wg.Done()
			}
		}(wg)
	}
	wg.Wait()
}

func hashFunc(data []byte, hf func() hash.Hash32, b *testing.B) {
	for i := 0; i < b.N; i++ {
		h := hf()
		h.Write(data)
		h.Sum32()
	}
}

func BenchmarkAdler(b *testing.B) {
	hashFunc(s, adler32.New, b)
}

func BenchmarkFnv(b *testing.B) {
	hashFunc(s, fnv.New32, b)
}

func BenchmarkAdlerLong(b *testing.B) {
	hashFunc(l, adler32.New, b)
}

func BenchmarkFnvLong(b *testing.B) {
	hashFunc(l, fnv.New32, b)
}
