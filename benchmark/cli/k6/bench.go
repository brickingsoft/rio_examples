package k6

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
)

func Bench(host string, port int, count int, dur string, out string) {
	if host == "" || port < 1 {
		fmt.Println("k6: host and port are required")
		return
	}

	rates := make(map[string]float64)
	buf := new(bytes.Buffer)
	for _, server := range srv.HttpServers {
		port++

		buf.WriteString("------" + server.Name + "------\n")
		scriptFile := ScriptFile(server.Name, host, port, count, dur)
		// k6 run file.js
		cmd := exec.Command("k6", "run", "-q", scriptFile)
		outBuf := new(bytes.Buffer)
		cmd.Stderr = outBuf
		cmd.Stdout = outBuf
		fmt.Println("["+server.Name+"]", "[HTTP]", "[BED]")
		err := cmd.Run()
		txt := outBuf.Bytes()
		fmt.Println(string(txt))
		fmt.Println("["+server.Name+"]", "[HTTP]", "[END]")
		_ = os.Remove(scriptFile)
		if err != nil {
			fmt.Println("k6 failed, check it is installed")
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
	textFile := filepath.Join(out, "benchmark_k6.txt")
	_ = os.WriteFile(textFile, buf.Bytes(), 0644)
	// write image
	req := images.Request{
		Path:  filepath.Join(out, "benchmark_k6.png"),
		Title: fmt.Sprintf("Benchmark(%s)", "k6"),
		Label: "r/s",
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

func parseRate(b []byte) (v float64) {
	buf := bytes.NewBuffer(b)
	reader := bufio.NewReader(buf)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			return
		}
		ls := strings.TrimSpace(string(line))
		if strings.HasPrefix(ls, "http_reqs") {
			rate, rateErr := strconv.ParseFloat(strings.Split(strings.Split(ls, ": ")[1], " ")[0], 64)
			if rateErr != nil {
				fmt.Println("parse rate failed:", rateErr)
			}
			return rate
		}
	}
}
