package zk

type Config struct {
	Host           []string
	Auth           string
	SessionTimeout float64
	RootPath       string
}
