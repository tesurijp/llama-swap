package main

import (
    _ "embed"
    "os"
    "os/exec"
    "path/filepath"
    "syscall"
    "flag"

    "github.com/getlantern/systray"
    "golang.org/x/sys/windows"
)

const (
    TargetProgram = "llama-swap.exe"
    TargetURL = "http://localhost"
)

var (
    childProcess *os.Process
    listenStr *string
    //go:embed icon.ico
    iconData []byte
)


func main() {
    listenStr = flag.String("listen", ":8080", "listen ip/port")
    flag.Parse()

    err := runTargetProgram();
    if( err == nil ) {
        systray.Run(onReady, onExit)
    }
}

func onReady() {
    if childProcess == nil {
        systray.Quit()
        return
    }

    systray.SetIcon(iconData)
    systray.SetTitle("llamaSwap")
    systray.SetTooltip("llamaSwap")

    mOpenLog := systray.AddMenuItem("Open log", "Open llamaSwap-Web log page")
    mOpenModel := systray.AddMenuItem("Open model", "Open llamaSwap-Web model page")
    mTerminateChild := systray.AddMenuItem("Exit", "Ext llamaSwap")

    go func() {
        for {
            select {
            case <-mOpenLog.ClickedCh:
                openBrowser("/ui")

            case <-mOpenModel.ClickedCh:
                openBrowser("/ui/models")

            case <-mTerminateChild.ClickedCh:
                terminateChildProcess()
            }
        }
    }()

    go func() {
        if childProcess != nil {
            childProcess.Wait()
            systray.Quit()
        }
    }()
}

func onExit() {
    terminateChildProcess()
}

func runTargetProgram() error {
    programPath, err := filepath.Abs(TargetProgram)
    if(err == nil) {
        cmd := exec.Command(programPath,os.Args[1:] ...)
        cmd.SysProcAttr = &syscall.SysProcAttr{
            HideWindow:    true,
            CreationFlags: windows.CREATE_NO_WINDOW,
        }
        err = cmd.Start()
        childProcess = cmd.Process
    }

    return err 
}

func openBrowser(page string) {
    exec.Command("rundll32", "url.dll,FileProtocolHandler", TargetURL+ *listenStr + page ).Start()
}

func terminateChildProcess() {
    if childProcess != nil {
        pid := childProcess.Pid
        p, err := os.FindProcess(pid)
        if err == nil {
            p.Kill()
        }
        childProcess = nil
    }
}

