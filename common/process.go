package common

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

func RunCmd(ctx context.Context, cmd string, args ...string) ([]byte, []byte, error) {
	cmd = fmt.Sprintf("%s %s", cmd, strings.Join(args, " "))
	//ins := exec.CommandContext(ctx, "bash", "-c", cmd)	//如果启动了子进程，则ctx cancel无法关闭子进程
	ins := exec.Command("bash", "-c", cmd)

	//子进程设定到新进程组（与孙/孙...进程并入到一个组，方便后面统一kill）
	ins.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	exitCh := make(chan struct{}, 1)
	defer func() {
		exitCh <- struct{}{}
		close(exitCh)
	}()
	go func() {
		for {
			select {
			case <-exitCh:
				return
			case <-ctx.Done():
				if ins.Process == nil {
					//等待子进程创建
					time.Sleep(10 * time.Millisecond)
				} else {
					syscall.Kill(-ins.Process.Pid, syscall.SIGKILL)
					return
				}
			}
		}
	}()

	//后面优化成对象池
	var stdout, stderr bytes.Buffer
	ins.Stdout = &stdout
	ins.Stderr = &stderr

	if err := ins.Run(); err != nil {
		return stdout.Bytes(), stderr.Bytes(), errors.Wrap(err, "run cmd")
	}

	return stdout.Bytes(), stderr.Bytes(), nil
}
