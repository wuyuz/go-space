package main

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"strings"
	"testing"
	"time"

	"gopkg.in/tomb.v1"
)

func WithEchoServer(t *testing.T, f func(string, chan []byte)) {
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal("Failed to create TCP server", err)
	}

	defer ln.Close()

	response := make(chan []byte, 1)
	tomb := tomb.Tomb{}

	go func() {
		defer tomb.Done()
		src, err := ln.Accept()
		if err != nil {
			select {
			case <-tomb.Dying():
			default:
				t.Fatal("Failed to accept client")
			}
			return
		}

		ln.Close()

		scan := bufio.NewScanner(src)
		if scan.Scan() {
			received := append(scan.Bytes(), '\n')
			response <- received

			src.Write(received)
		}
	}()

	f(ln.Addr().String(), response)

	tomb.Killf("Function body finished")
	ln.Close()
	tomb.Wait()

	close(response)
}

func WithEchoProxy(t *testing.T, f func(proxy net.Conn, response chan []byte, proxyServer *Proxy)) {
	WithEchoServer(t, func(upstream string, response chan []byte) {
		proxy := NewTestProxy("test", upstream)
		proxy.Start()

		conn, err := net.Dial("tcp", proxy.Listen)
		if err != nil {
			t.Error("Unable to dial TCP server", err)
		}

		f(conn, response, proxy)

		proxy.Stop()
	})
}

func AssertDeltaTime(t *testing.T, message string, actual, expected, delta time.Duration) {
	diff := actual - expected
	if diff < 0 {
		diff *= -1
	}
	if diff > delta {
		t.Errorf("[%s] Time was more than %v off: got %v expected %v", message, delta, actual, expected)
	} else {
		t.Logf("[%s] Time was correct: %v (expected %v)", message, actual, expected)
	}
}

func DoLatencyTest(t *testing.T, upLatency, downLatency *LatencyToxic) {
	WithEchoProxy(t, func(conn net.Conn, response chan []byte, proxy *Proxy) {
		t.Logf("Using latency: Up: %dms +/- %dms, Down: %dms +/- %dms", upLatency.Latency, upLatency.Jitter, downLatency.Latency, downLatency.Jitter)
		proxy.upToxics.SetToxicValue(upLatency)
		proxy.downToxics.SetToxicValue(downLatency)

		msg := []byte("hello world " + strings.Repeat("a", 32*1024) + "\n")

		timer := time.Now()
		_, err := conn.Write(msg)
		if err != nil {
			t.Error("Failed writing to TCP server", err)
		}

		resp := <-response
		if !bytes.Equal(resp, msg) {
			t.Error("Server didn't read correct bytes from client:", string(resp))
		}
		AssertDeltaTime(t,
			"Server read",
			time.Now().Sub(timer),
			time.Duration(upLatency.Latency)*time.Millisecond,
			time.Duration(upLatency.Jitter+10)*time.Millisecond,
		)
		timer2 := time.Now()

		scan := bufio.NewScanner(conn)
		if scan.Scan() {
			resp = append(scan.Bytes(), '\n')
			if !bytes.Equal(resp, msg) {
				t.Error("Client didn't read correct bytes from server:", string(resp))
			}
		}
		AssertDeltaTime(t,
			"Client read",
			time.Now().Sub(timer2),
			time.Duration(downLatency.Latency)*time.Millisecond,
			time.Duration(downLatency.Jitter+10)*time.Millisecond,
		)
		AssertDeltaTime(t,
			"Round trip",
			time.Now().Sub(timer),
			time.Duration(upLatency.Latency+downLatency.Latency)*time.Millisecond,
			time.Duration(upLatency.Jitter+downLatency.Jitter+10)*time.Millisecond,
		)

		upLatency.Enabled = false
		downLatency.Enabled = false
		proxy.upToxics.SetToxicValue(upLatency)
		proxy.downToxics.SetToxicValue(downLatency)

		err = conn.Close()
		if err != nil {
			t.Error("Failed to close TCP connection", err)
		}
	})
}

func TestUpstreamLatency(t *testing.T) {
	DoLatencyTest(t, &LatencyToxic{Enabled: true, Latency: 100}, &LatencyToxic{Enabled: false})
}

func TestDownstreamLatency(t *testing.T) {
	DoLatencyTest(t, &LatencyToxic{Enabled: false}, &LatencyToxic{Enabled: true, Latency: 100})
}

func TestFullstreamLatencyEven(t *testing.T) {
	DoLatencyTest(t, &LatencyToxic{Enabled: true, Latency: 100}, &LatencyToxic{Enabled: true, Latency: 100})
}

func TestFullstreamLatencyBiasUp(t *testing.T) {
	DoLatencyTest(t, &LatencyToxic{Enabled: true, Latency: 1000}, &LatencyToxic{Enabled: true, Latency: 100})
}

func TestFullstreamLatencyBiasDown(t *testing.T) {
	DoLatencyTest(t, &LatencyToxic{Enabled: true, Latency: 100}, &LatencyToxic{Enabled: true, Latency: 1000})
}

func TestZeroLatency(t *testing.T) {
	DoLatencyTest(t, &LatencyToxic{Enabled: true, Latency: 0}, &LatencyToxic{Enabled: true, Latency: 0})
}

func AssertEchoResponse(t *testing.T, client, server net.Conn) {
	msg := []byte("hello world\n")

	_, err := client.Write(msg)
	if err != nil {
		t.Error("Failed writing to TCP server", err)
	}

	scan := bufio.NewScanner(server)
	if !scan.Scan() {
		t.Error("Client unexpectedly closed connection")
	}

	resp := append(scan.Bytes(), '\n')
	if !bytes.Equal(resp, msg) {
		t.Error("Server didn't read correct bytes from client:", string(resp))
	}

	_, err = server.Write(resp)
	if err != nil {
		t.Error("Failed writing to TCP client", err)
	}

	scan = bufio.NewScanner(client)
	if !scan.Scan() {
		t.Error("Server unexpectedly closed connection")
	}

	resp = append(scan.Bytes(), '\n')
	if !bytes.Equal(resp, msg) {
		t.Error("Client didn't read correct bytes from server:", string(resp))
	}
}

func TestPersistentConnections(t *testing.T) {
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal("Failed to create TCP server", err)
	}

	defer ln.Close()

	proxy := NewTestProxy("test", ln.Addr().String())
	proxy.Start()
	defer proxy.Stop()

	serverConnRecv := make(chan net.Conn)

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			t.Error("Unable to accept TCP connection", err)
		}
		serverConnRecv <- conn
	}()

	conn, err := net.Dial("tcp", proxy.Listen)
	if err != nil {
		t.Error("Unable to dial TCP server", err)
	}

	serverConn := <-serverConnRecv

	proxy.upToxics.SetToxicValue(&LatencyToxic{Enabled: true, Latency: 0})
	proxy.downToxics.SetToxicValue(&LatencyToxic{Enabled: true, Latency: 0})

	AssertEchoResponse(t, conn, serverConn)

	proxy.upToxics.ResetToxics()
	proxy.downToxics.ResetToxics()

	AssertEchoResponse(t, conn, serverConn)

	proxy.upToxics.ResetToxics()
	proxy.downToxics.ResetToxics()

	AssertEchoResponse(t, conn, serverConn)

	err = conn.Close()
	if err != nil {
		t.Error("Failed to close TCP connection", err)
	}
}

func TestLatencyToxicCloseRace(t *testing.T) {
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal("Failed to create TCP server", err)
	}

	defer ln.Close()

	proxy := NewTestProxy("test", ln.Addr().String())
	proxy.Start()
	defer proxy.Stop()

	go func() {
		for {
			_, err := ln.Accept()
			if err != nil {
				return
			}
		}
	}()

	// Check for potential race conditions when interrupting toxics
	for i := 0; i < 1000; i++ {
		proxy.upToxics.SetToxicValue(&LatencyToxic{Enabled: true, Latency: 10})
		conn, err := net.Dial("tcp", proxy.Listen)
		if err != nil {
			t.Error("Unable to dial TCP server", err)
		}
		conn.Write([]byte("hello"))
		conn.Close()
		proxy.upToxics.SetToxicValue(&LatencyToxic{Enabled: false})
	}
}

func TestLatencyToxicBandwidth(t *testing.T) {
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal("Failed to create TCP server", err)
	}

	defer ln.Close()

	proxy := NewTestProxy("test", ln.Addr().String())
	proxy.Start()
	defer proxy.Stop()

	buf := []byte(strings.Repeat("hello world ", 1000))

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			t.Error("Unable to accept TCP connection", err)
		}
		for err == nil {
			_, err = conn.Write(buf)
		}
	}()

	conn, err := net.Dial("tcp", proxy.Listen)
	if err != nil {
		t.Error("Unable to dial TCP server", err)
	}

	proxy.downToxics.SetToxicValue(&LatencyToxic{Enabled: true, Latency: 100})

	time.Sleep(100 * time.Millisecond) // Wait for latency toxic
	buf2 := make([]byte, len(buf))

	start := time.Now()
	count := 0
	for i := 0; i < 100; i++ {
		n, err := io.ReadFull(conn, buf2)
		count += n
		if err != nil {
			t.Error(err)
			break
		}
	}

	// Assert the transfer was at least 100MB/s
	AssertDeltaTime(t, "Latency toxic bandwidth", time.Since(start), 0, time.Duration(count/100000)*time.Millisecond)

	err = conn.Close()
	if err != nil {
		t.Error("Failed to close TCP connection", err)
	}
}

func TestProxyLatency(t *testing.T) {
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal("Failed to create TCP server", err)
	}

	defer ln.Close()

	proxy := NewTestProxy("test", ln.Addr().String())
	proxy.Start()
	defer proxy.Stop()

	serverConnRecv := make(chan net.Conn)

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			t.Error("Unable to accept TCP connection", err)
		}
		serverConnRecv <- conn
	}()

	conn, err := net.Dial("tcp", proxy.Listen)
	if err != nil {
		t.Error("Unable to dial TCP server", err)
	}

	serverConn := <-serverConnRecv

	start := time.Now()
	for i := 0; i < 100; i++ {
		AssertEchoResponse(t, conn, serverConn)
	}
	latency := time.Now().Sub(start) / 200
	if latency > 300*time.Microsecond {
		t.Errorf("Average proxy latency > 300µs (%v)", latency)
	} else {
		t.Logf("Average proxy latency: %v", latency)
	}

	err = conn.Close()
	if err != nil {
		t.Error("Failed to close TCP connection", err)
	}
}

func TestBandwidthToxic(t *testing.T) {
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal("Failed to create TCP server", err)
	}

	defer ln.Close()

	proxy := NewTestProxy("test", ln.Addr().String())
	proxy.Start()
	defer proxy.Stop()

	serverConnRecv := make(chan net.Conn)

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			t.Error("Unable to accept TCP connection", err)
		}
		serverConnRecv <- conn
	}()

	conn, err := net.Dial("tcp", proxy.Listen)
	if err != nil {
		t.Error("Unable to dial TCP server", err)
	}

	serverConn := <-serverConnRecv

	rate := 1000 // 1MB/s
	proxy.upToxics.SetToxicValue(&BandwidthToxic{Enabled: true, Rate: int64(rate)})

	buf := []byte(strings.Repeat("hello world ", 40000)) // 480KB
	go func() {
		n, err := conn.Write(buf)
		conn.Close()
		if n != len(buf) || err != nil {
			t.Errorf("Failed to write buffer: (%d == %d) %v", n, len(buf), err)
		}
	}()

	buf2 := make([]byte, len(buf))
	start := time.Now()
	_, err = io.ReadAtLeast(serverConn, buf2, len(buf2))
	if err != nil {
		t.Errorf("Proxy read failed: %v", err)
	} else if bytes.Compare(buf, buf2) != 0 {
		t.Errorf("Server did not read correct buffer from client!")
	}

	AssertDeltaTime(t,
		"Bandwidth",
		time.Now().Sub(start),
		time.Duration(len(buf))*time.Second/time.Duration(rate*1000),
		10*time.Millisecond,
	)
}

func TestSlicerToxic(t *testing.T) {
	data := []byte(strings.Repeat("hello world ", 40000)) // 480 kb
	slicer := &SlicerToxic{Enabled: true, AverageSize: 1024, SizeVariation: 512, Delay: 10}

	input := make(chan *StreamChunk)
	output := make(chan *StreamChunk)
	stub := NewToxicStub(input, output)

	done := make(chan bool)
	go func() {
		slicer.Pipe(stub)
		done <- true
	}()
	defer func() {
		input <- nil
		<-done
	}()

	input <- &StreamChunk{data: data}

	buf := make([]byte, 0, len(data))
	reads := 0
L:
	for {
		select {
		case c := <-output:
			reads++
			buf = append(buf, c.data...)
		case <-time.After(5 * time.Millisecond):
			break L
		}
	}

	if reads < 480/2 || reads > 480/2+480 {
		t.Errorf("Expected to read about 480 times, but read %d times.", reads)
	}
	if bytes.Compare(buf, data) != 0 {
		t.Errorf("Server did not read correct buffer from client!")
	}
}

func TestToxicUpdate(t *testing.T) {
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal("Failed to create TCP server", err)
	}

	defer ln.Close()

	proxy := NewTestProxy("test", ln.Addr().String())
	proxy.Start()
	defer proxy.Stop()

	serverConnRecv := make(chan net.Conn)

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			t.Error("Unable to accept TCP connection", err)
		}
		serverConnRecv <- conn
	}()

	conn, err := net.Dial("tcp", proxy.Listen)
	if err != nil {
		t.Error("Unable to dial TCP server", err)
	}

	serverConn := <-serverConnRecv

	running := make(chan struct{})
	go func() {
		enabled := false
		for {
			select {
			case <-running:
				return
			default:
				proxy.upToxics.SetToxicValue(&LatencyToxic{Enabled: enabled})
				enabled = !enabled
				proxy.downToxics.SetToxicValue(&LatencyToxic{Enabled: enabled})
			}
		}
	}()

	for i := 0; i < 100; i++ {
		AssertEchoResponse(t, conn, serverConn)
	}
	close(running)

	err = conn.Close()
	if err != nil {
		t.Error("Failed to close TCP connection", err)
	}
}

func BenchmarkBandwidthToxic100MB(b *testing.B) {
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		b.Fatal("Failed to create TCP server", err)
	}

	defer ln.Close()

	proxy := NewTestProxy("test", ln.Addr().String())
	proxy.Start()
	defer proxy.Stop()

	buf := []byte(strings.Repeat("hello world ", 1000))

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			b.Error("Unable to accept TCP connection", err)
		}
		buf2 := make([]byte, len(buf))
		for err == nil {
			_, err = conn.Read(buf2)
		}
	}()

	conn, err := net.Dial("tcp", proxy.Listen)
	if err != nil {
		b.Error("Unable to dial TCP server", err)
	}

	proxy.upToxics.SetToxicValue(&BandwidthToxic{Enabled: true, Rate: 100 * 1000})

	b.SetBytes(int64(len(buf)))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n, err := conn.Write(buf)
		if err != nil || n != len(buf) {
			b.Errorf("%v, %d == %d", err, n, len(buf))
			break
		}
	}

	err = conn.Close()
	if err != nil {
		b.Error("Failed to close TCP connection", err)
	}
}

func BenchmarkNoopToxic(b *testing.B) {
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		b.Fatal("Failed to create TCP server", err)
	}

	defer ln.Close()

	proxy := NewTestProxy("test", ln.Addr().String())
	proxy.Start()
	defer proxy.Stop()

	buf := []byte(strings.Repeat("hello world ", 1000))

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			b.Error("Unable to accept TCP connection", err)
		}
		buf2 := make([]byte, len(buf))
		for err == nil {
			_, err = conn.Read(buf2)
		}
	}()

	conn, err := net.Dial("tcp", proxy.Listen)
	if err != nil {
		b.Error("Unable to dial TCP server", err)
	}

	b.SetBytes(int64(len(buf)))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n, err := conn.Write(buf)
		if err != nil || n != len(buf) {
			b.Errorf("%v, %d == %d", err, n, len(buf))
			break
		}
	}

	err = conn.Close()
	if err != nil {
		b.Error("Failed to close TCP connection", err)
	}
}
