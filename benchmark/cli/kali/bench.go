package kali

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/brickingsoft/rio_examples/benchmark/images"
	"github.com/brickingsoft/rio_examples/benchmark/srv"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func Bench(host string, port int, count int, repeat int, dur string, out string) {
	if host == "" || port < 1 {
		fmt.Println("tcpkali: host and port are required")
		return
	}

	kind := ""

	var d time.Duration
	var dErr error
	if repeat == 0 {
		dur = strings.ToLower(strings.TrimSpace(dur))
		if dur == "" {
			dur = "10s"
		}
		d, dErr = time.ParseDuration(dur)
		if dErr != nil {
			fmt.Println("parse time failed:", dErr)
			fmt.Println("use 10s.")
			d = 10 * time.Second
		}
		kind = fmt.Sprintf("C%dT%s", count, d.String())
	} else {
		kind = fmt.Sprintf("C%dR%s", count, formatNumber(repeat))
	}

	rates := make(map[string]float64)

	buf := new(bytes.Buffer)
	for _, server := range srv.TcpServers {
		port++

		buf.WriteString("------" + server.Name + "------\n")

		var cmd *exec.Cmd
		if repeat > 1 {
			// tcpkali --workers 1 -c 50 -r 5k -m "PING"  192.168.100.1:9000
			cmd = exec.Command("tcpkali",
				"--workers", "1",
				"-c", strconv.Itoa(count),
				"-r", strconv.Itoa(repeat),
				"-m", fmt.Sprintf("\"%s\"", "PING"),
				fmt.Sprintf("%s:%d", host, port),
			)
		} else {
			// tcpkali --workers 1 -c 50 -T 10s -m "PING"  192.168.100.1:9000
			cmd = exec.Command("tcpkali",
				"--workers", "1",
				"-c", strconv.Itoa(count),
				"-T", d.String(),
				"-m", fmt.Sprintf("\"%s\"", "PING"),
				fmt.Sprintf("%s:%d", host, port),
			)
		}

		outBuf := new(bytes.Buffer)
		cmd.Stderr = outBuf
		cmd.Stdout = outBuf
		fmt.Println("["+server.Name+"]", "["+kind+"]", "[BED]")
		err := cmd.Run()
		txt := outBuf.Bytes()
		fmt.Println(string(txt))
		fmt.Println("["+server.Name+"]", "["+kind+"]", "[END]")
		if err != nil {
			fmt.Println("tcpkali failed, check it is installed")
			fmt.Println(fmt.Sprintf("%+v", err))
			return
		}
		var rate float64
		if outBuf.Len() > 0 {
			buf.Write(txt)
			buf.WriteString("\n")
			rate = parseRate(txt)
		}
		rates[server.Name] = rate
	}

	_ = os.MkdirAll(filepath.Dir(out), 0644)
	// write text
	textFile := filepath.Join(out, fmt.Sprintf("benchmark_tcpkali_%s.txt", kind))
	_ = os.WriteFile(textFile, buf.Bytes(), 0644)
	// write image
	req := images.Request{
		Path:  filepath.Join(out, fmt.Sprintf("benchmark_tcpkali_%s.png", kind)),
		Title: fmt.Sprintf("Benchmark(%s)", kind),
		Label: "pps",
		Items: make([]images.Item, 0, 1),
	}
	for title, rate := range rates {
		req.Items = append(req.Items, images.Item{
			Name:  title,
			Value: rate,
		})
	}

	err := images.Draw(req)
	if err != nil {
		fmt.Println(err)
	}

}

func formatNumber(n int) string {
	unit := ""
	value := float64(n)

	switch {
	case n >= 1000:
		unit = "K"
		value = value / 1000
		break
	default:
		break
	}

	result := strconv.FormatFloat(value, 'f', 1, 64)
	result = strings.TrimSuffix(result, ".0")
	return result + unit
}

func parseRate(b []byte) (v float64) {
	buf := bytes.NewBuffer(b)
	reader := bufio.NewReader(buf)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			return
		}
		ls := strings.TrimSpace(string(line))
		if strings.HasPrefix(ls, "Packet rate estimate:") {
			rate, rateErr := strconv.ParseFloat(strings.Split(strings.Split(ls, ": ")[1], "↓,")[0], 64)
			if rateErr != nil {
				fmt.Println("parse rate failed:", rateErr)
			}
			return rate
		}
	}
}
