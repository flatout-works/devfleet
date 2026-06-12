package controller

import "strings"

func kataPreflightCommand(cmd []string) []string {
	script := strings.Join([]string{
		"set -u",
		`echo "[chetter] preflight: pwd=$(pwd)"`,
		`echo "[chetter] preflight: workspace"`,
		`ls -la /workspace || true`,
		`echo "[chetter] preflight: resolv.conf"`,
		`cat /etc/resolv.conf || true`,
		`echo "[chetter] preflight: proxy=$CHETTER_PROXY"`,
		`curl -sS --connect-timeout 5 -x "http://$CHETTER_PROXY" -o /dev/null -w "[chetter] preflight: synthetic http=%{http_code}\n" https://api.synthetic.new/openai/v1/models || echo "[chetter] preflight: synthetic proxy failed"`,
		`curl --noproxy '*' -sS --connect-timeout 5 -o /dev/null -w "[chetter] preflight: synthetic direct http=%{http_code}\n" https://api.synthetic.new/openai/v1/models || echo "[chetter] preflight: synthetic direct failed"`,
		`curl --noproxy '*' -sS --connect-timeout 5 --max-time 30 -o /tmp/chetter-chat.json -w "[chetter] preflight: chat http=%{http_code}\n" -H "Authorization: Bearer $SYNTHETIC_API_KEY" -H "Content-Type: application/json" -d '{"model":"hf:zai-org/GLM-5.1","messages":[{"role":"user","content":"Reply with OK only."}],"max_tokens":16}' https://api.synthetic.new/openai/v1/chat/completions || echo "[chetter] preflight: chat failed"`,
		`test -s /tmp/chetter-chat.json && head -c 500 /tmp/chetter-chat.json && echo || true`,
		`echo "[chetter] preflight: opencode version"`,
		`opencode --version 2>&1 || echo "[chetter] preflight: opencode version failed"`,
		`echo "[chetter] preflight: TCP connectivity (bash) to api.synthetic.new:443"`,
		`timeout 10 bash -c 'exec 3<>/dev/tcp/api.synthetic.new/443 2>/dev/null && echo BASH_TCP_OK || echo BASH_TCP_FAIL' 2>&1 || echo "[chetter] preflight: bash tcp timeout=$?"`,
		`echo "[chetter] preflight: IPv4-only curl test"`,
		`curl -4 --noproxy '*' -sS --connect-timeout 5 -o /dev/null -w "[chetter] preflight: curl4 http=%{http_code}\n" https://api.synthetic.new/openai/v1/models || echo "[chetter] preflight: curl4 failed"`,
		`echo "[chetter] preflight: IPv6-only curl test"`,
		`curl -6 --noproxy '*' -sS --connect-timeout 5 -o /dev/null -w "[chetter] preflight: curl6 http=%{http_code}\n" https://api.synthetic.new/openai/v1/models || echo "[chetter] preflight: curl6 failed (expected in Kata)"`,
		`echo "[chetter] preflight: exec opencode from /tmp (tmpfs) with --dir /workspace (timeout 300s + strace)"`,
		"cd /tmp && strace -f -e trace=network,read,write,process,select,poll,epoll_wait,connect,accept,getpeername -tt -y -s 128 -o /host-tmp/strace.log timeout 300 " + shellQuoteArgs(cmd) + " < /dev/null || echo '[chetter] preflight: opencode timed out or exited non-zero'",
		`echo "[chetter] preflight: strace summary:"`,
		`grep -c "+++ exited" /host-tmp/strace.log 2>/dev/null || echo "no strace log"`,
		`tail -n 5 /host-tmp/strace.log 2>/dev/null || true`,
	}, "\n")
	return []string{"sh", "-c", script}
}
