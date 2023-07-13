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

type RunCmdResult struct {
	StdOut []byte `json:"std_out"`
	StdErr []byte `json:"std_err"`
}

func RunCmd(ctx context.Context, cmd string, args ...string) (*RunCmdResult, error) {
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
		return nil, errors.Wrap(err, "run cmd")
	}

	var result RunCmdResult
	result.StdOut = stdout.Bytes()
	result.StdErr = stderr.Bytes()

	return &result, nil
}

//apple script
func AppleScript(ctx context.Context, script string, isJxa bool) (*RunCmdResult, error) {
	args := []string{}
	if isJxa {
		args = append(args, "-l", "JavaScript")
	}
	args = append(args, "-e", script)
	result, err := RunCmd(ctx, "osascript", args...)

	return result, errors.Wrap(err, "run apple script")
}

//convinient

//创建提醒
func MakeRmd(ctx context.Context, title, body string, minute int) (*RunCmdResult, error) {
	script := fmt.Sprintf(`tell application "Reminders"
	set theReminder to make new reminder with properties {name:"%s", body:"%s", remind me date:(current date) + (%d * minutes)}
	end tell`, title, body, minute)
	return AppleScript(ctx, script, false)
}

//发送系统通知
func MakeNotify(ctx context.Context, title, content string) (*RunCmdResult, error) {
	script := fmt.Sprintf(`display notification "%s" with title "%s"`, content, title)
	return AppleScript(ctx, script, false)
}
