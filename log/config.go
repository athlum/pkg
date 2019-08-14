package log

type Config struct {
	EndPoint string
	Protocol string
	LogFile  string
	Verbose  int
	Debug    bool
	AppId    string
	Disable  bool
}
