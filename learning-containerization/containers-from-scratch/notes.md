### SysProcAttr

- A new child process in Go can be spawned with a following commands:

```go
cmd := exec.Command(os.Args[2], os.Args[3]...)
cmd.Stdin = os.Stdin
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr

cmd.SysProcAttr = &syscall.SysProcAttr{
    Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,}

cmd.Run()
```
- `cmd` represents a struct from `os/exec` package, which represents the external command to be run.

    - `cmd.Stdin/Stdout/Stderr`: Write up the I/O terminal. Connects the containers input/output and errors directly to the host terminal.
    - `cmd.SysProcAttr`: struct that can be used to modify OS Level process attributes.
    - `cmd.Env` : Environment variable for the process.
    - `cmd.Dir` : Working directory for the process (?).

- `cmd` have certain methods to control it:
    
    - `cmd.Run()`: Start and Wait for the process/command to finish.
    - `cmd.Start()`: Start with out waiting.
    - `cmd.Wait()`: Wait after `Start()`.

- Clone flags tells what namespaces to create for the new process.

- The above codes only spawns a child process in the new namespaces, that we create through clone flags.
    - However, if we require to customize the namespaces, such as setting a custom hostname, and then run the process, a child process is required.
    - The child process will perform all the actions to configure namespaces.
    - The common way to connect a paraent and child process is via **symlink** `/proc/self/exe`.

    - `/proc/self/exe` always points to the executable file of the currently running process.
    

```go
func parent(){

    cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
    cmd.Stdin = os.Stdin
    .
    .

    // isolated namspaces for the child
    cmd.SysProcAttr = &syscall.SysProcAttr{
    Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
    }

    cmd.Run()
}

func child(){

    syscall.Sethostname([]byte('my-container'))

    cmd := exec.Command(os.Args[2], os.Args[3:]...)
    cmd.Stdin = os.Stdin
    .
    .
    .

    cmd.Run()
}

```










