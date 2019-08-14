package utils

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"math"
	"net"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"sync/atomic"
	"sync"
	"time"
)

var (
	ERROR_MissingConfigFile   = errors.New("unknown config file path.")
	ERROR_ConfigFileNotExists = errors.New("config file does not exist.")
)

func IntAbs(n int) int {
	return int(math.Abs(float64(n)))
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

var localIp = &atomic.Value{}
var localIpCache = &sync.Once{}
const defaultLocalIpDuration = time.Second * 30

func flushLocalIp(d time.Duration) {
	t := time.NewTicker(d)
	defer t.Stop()
	for {
		select{
		case <- t.C:
			localIp.Store(GetLocalIP())
		}
	}
}

func GetLocalIPCached(ds ...time.Duration) string {
	localIpCache.Do(func() {
		localIp.Store(GetLocalIP())
		d := defaultLocalIpDuration
		if len(ds) > 0 && ds[0] > 0 {
			d = ds[0]
		}
		go flushLocalIp(d)
	})
	return localIp.Load().(string)
}

func ShortSha1(data []byte) []byte {
	d := sha1.Sum(data)
	return d[:6]
}

func Duration(seconds float64) time.Duration {
	return time.Millisecond * time.Duration(seconds*1000)
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func ParseIP(address string) (net.IP, error) {
	addr, err := net.LookupIP(address)
	if err != nil {
		return []byte{}, err
	}
	if len(addr) < 1 {
		return []byte{}, fmt.Errorf("failed to parse IP from address '%v'", address)
	}
	return addr[0], nil
}

func FloatToString(input_num float64) string {
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func ParseTime(timeStr string) (time.Time, error) {
	timeStr = strings.Replace(timeStr, " ", "T", -1)
	timeStr = fmt.Sprintf("%s+08:00", timeStr)
	return time.Parse(time.RFC3339, timeStr)
}

func GenerateUUID() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("failed to read random bytes: %v", err)
	}
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%12x",
		buf[0:4],
		buf[4:6],
		buf[6:8],
		buf[8:10],
		buf[10:16],
	), nil
}

func JoinUrl(src, dest string) (string, error) {
	u, err := url.Parse(src)
	if err != nil {
		return "", err
	}
	u.Path = path.Join(u.Path, dest)
	return u.String(), nil
}

func NextBackOff(backoff, factor, unit, max float64, backoffCount int) (float64, int) {
	backoffCount += 1
	nextBackOff := backoff + float64(backoffCount)*factor*unit
	if nextBackOff > max {
		return max, backoffCount
	}
	return nextBackOff, backoffCount
}

func OnlyPrefix(str, prefix string) bool {
	return strings.HasPrefix(str, prefix) && !strings.HasSuffix(str, prefix)
}

func SplitCommand(command string) []string {
	cmds := strings.Split(command, " ")
	splited := []string{}
	block := ""
	quote := ""
	for _, c := range cmds {
		if quote == "" {
			if OnlyPrefix(c, "'") || OnlyPrefix(c, "\"") {
				quote = string(c[0])
				block = strings.TrimPrefix(c, quote) + " "
			} else {
				splited = append(splited, c)
			}
		} else {
			if strings.HasSuffix(c, quote) {
				block += strings.TrimSuffix(c, quote)
				splited = append(splited, block)
				block = ""
				quote = ""
			} else {
				block += c + " "
			}
		}
	}
	return splited
}

func MetricDuration(start time.Time) int {
	return int(time.Now().Sub(start) / time.Microsecond)
}

func IfNil(v, d string) string {
	if v == "" {
		return d
	}
	return v
}

func ReadFile(fp, env string) ([]byte, error) {
	if fp == "" && env != "" {
		fp = os.Getenv(env)
	}
	if fp == "" {
		return nil, ERROR_MissingConfigFile
	}
	if _, err := os.Stat(fp); err != nil {
		return nil, ERROR_ConfigFileNotExists
	}

	data, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config file:")
	}
	return data, nil
}
