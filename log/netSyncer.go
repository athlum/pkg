package log

import (
	"go.uber.org/zap/zapcore"
	"net"
)

const (
	TCP = "tcp"
	UDP = "udp"
)

type AsyncNetWriter struct {
	endpoint string
	protocol string
}

func NewAsyncNetWriter(endpoint, protocol string) *AsyncNetWriter {
	return &AsyncNetWriter{endpoint: endpoint, protocol: protocol}
}

func NewAsyncTcpWriter(endpoint string) *AsyncNetWriter {
	return NewAsyncNetWriter(endpoint, TCP)
}

func NewAsyncUdpWriter(endpoint string) *AsyncNetWriter {
	return NewAsyncNetWriter(endpoint, UDP)
}

func (anw *AsyncNetWriter) WriterSyncer() zapcore.WriteSyncer {
	return zapcore.AddSync(anw)
}

func (anw *AsyncNetWriter) Write(p []byte) (n int, err error) {
	tmp := make([]byte, len(p))
	copy(tmp, p)
	go func() {
		conn, err := net.Dial(anw.protocol, anw.endpoint)
		if err != nil {
			return
		}
		defer conn.Close()
		conn.Write(tmp)
	}()
	return len(tmp), nil
}
