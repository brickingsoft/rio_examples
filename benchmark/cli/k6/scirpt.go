package k6

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	scriptTemplate = `import http from 'k6/http';
export let options = {
  vus: %d, 
  duration: '%s', 
};
// default 默认函数
export default function () {
  // 标头
  let params = { headers: { 'Content-Type': 'text/plain' } };

  var res=http.get("http://%s:%d/",params)
}`
)

func Script(host string, port int, count int, dur string) string {
	return fmt.Sprintf(scriptTemplate, count, dur, host, port)
}

func ScriptFile(name string, host string, port int, count int, dur string) string {
	s := Script(host, port, count, dur)
	tmp := os.TempDir()
	tmpFile := filepath.Join(tmp, fmt.Sprintf("%s.js", name))
	_ = os.WriteFile(tmpFile, []byte(s), 0644)
	return tmpFile
}
