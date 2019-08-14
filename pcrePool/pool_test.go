package pcrePool

import (
	"fmt"
	"github.com/athlum/golang-pkg-pcre/src/pkg/pcre"
	"testing"
)

func init() {
	Init()
}

func pcrePick(p, s string) {
	re := pcre.MustCompile(p, 0)
	m := re.MatcherString(s, 0)
	m.NamedStringMap()
}

func poolPick(p, s string) {
	re, err := Compile(p)
	if err != nil {
		panic(err)
	}
	defer re.Collect()
	m := re.MatcherString(s, 0)
	m.NamedStringMap()
}

func BenchmarkLongPcre(b *testing.B) {
	for n := 0; n < b.N; n += 1 {
		pcrePick("{hostname: (?<hostname>.*), ip: (?<ip>.*), topic: (?<topic>.*)} (?<source_msg>.*)", `{hostname: adca-mesos-32.vm.elenet.me, ip: 10.101.64.117, topic: arch.appos_agent} {"error":"request: Post http://127.0.0.1:1988/metrics?key=docker: net/http: request canceled (Client.Timeout exceeded while awaiting headers)","indice":"appos.agent","level":"error","log_source":"appos-agent","msg":"Send stats","time":"2018-06-18T03:15:41+08:00"}`)
	}
}

func BenchmarkLongPool(b *testing.B) {
	for n := 0; n < b.N; n += 1 {
		poolPick("{hostname: (?<hostname>.*), ip: (?<ip>.*), topic: (?<topic>.*)} (?<source_msg>.*)", `{hostname: adca-mesos-32.vm.elenet.me, ip: 10.101.64.117, topic: arch.appos_agent} {"error":"request: Post http://127.0.0.1:1988/metrics?key=docker: net/http: request canceled (Client.Timeout exceeded while awaiting headers)","indice":"appos.agent","level":"error","log_source":"appos-agent","msg":"Send stats","time":"2018-06-18T03:15:41+08:00"}`)
	}
}

func TestPcreMatchString(t *testing.T) {
	re, err := Compile("{hostname: (?<hostname>.*), ip: (?<ip>.*), topic: (?<topic>.*)} (?<source_msg>.*)")
	if err != nil {
		t.Error(err)
	}
	defer re.Collect()
	m := re.MatcherString(`{hostname: adca-mesos-32.vm.elenet.me, ip: 10.101.64.117, topic: arch.appos_agent} {"error":"request: Post http://127.0.0.1:1988/metrics?key=docker: net/http: request canceled (Client.Timeout exceeded while awaiting headers)","indice":"appos.agent","level":"error","log_source":"appos-agent","msg":"Send stats","time":"2018-06-18T03:15:41+08:00"}`, 0)
	fmt.Println(m.NamedStringMap(), m.Matches())
}
