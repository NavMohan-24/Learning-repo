## Setting up a container in GO code.

### Setting up the environment

- Spining up a docker container

```bash
docker run -it --privileged --name my-container ubuntu bash
```
`--privileged` grants full root access (required to perform syscalls)

- Copying files from the device to container

```bash
docker cp path/to/source my-container:/path/to/destination
```

- Execution of the code
```bash
go run main.go run \bin\bash
```

### Spawing a Child process in Go

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

### Configuring Child Process/Container

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
### Setting up a new file system

- It requires a new file system which will be mounted as root of the child process.

    - To setup a minimal ubuntu root file system:
    ```bash
    debootstrap --arch amd64 focal /path/to/new/rootfs
    ```

    Here, `focal` implies Ubuntu 20.04. `debootstrap` is the tool that installs a ubuntu base system.
    
    -  minimal alphine linux rootfs:

    ```bash
    curl -o alpine.tar.gz https://dl-cdn.alpinelinux.org/alpine/v3.19/releases/x86_64/alpine-minirootfs-3.19.0-x86_64.tar.gz
    tar -xzvf alpine.tar.gz -C /path/to/new/rootfs
    ```
    The alphine rootfs use `sh` instead of `bash`.

- Once the new root filesystem is created, the root directory and working directory of child process needs to be changed.

```go
func child(){

	syscall.Sethostname([]byte("container"))
	syscall.Chroot("/path/to/new/rootfs")
	syscall.Chdir("/") // set the new root as the working directory.

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	.
    .

	cmd.Run() 

}
```
- Low level container runtime (`runc`) perfroms similar actions while spinning up a container.

    - The filesystem of image specified is unpacked and copied in to the host. The container is then `chroot`-ed in to that path. 

    - Container runtimes use more sophisticated _Layered Filesystem_. Its lowest layer is immutable and shared across all containers (For eg: linux base image will be in the lowest layer). The highest layer, would be specific to the container. 


### Mounting the `/proc`

<!-- - A container created only by cloning UTS and PID name spaces, will still shows process in the host while inspecting with `ps`.

    - If a docker container is the host and child container is inside it, then `ps` inside the child container would crash. -->

- Ideally containers should be isolated from the host and should only show the process within it.

<!-- - This can be done by mounting the `/proc` directory is a new mount namespace. -->


- At this point, running a `ps` command on host and container gives similar results as the following.

On host `ps aux` gives all the process running:

```bash
USER         PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND
root           1  0.0  0.0   5052  3584 pts/0    Ss   09:46   0:00 bash
root         595  0.0  0.0   5012  3584 pts/1    Ss   09:50   0:00 /bin/bash
root       31162  0.3  0.2 1273268 18440 pts/1   Sl   16:20   0:00 go run main.go run /bin/bash
root       31196  0.0  0.0 1227308 2176 pts/1    Sl   16:20   0:00 /tmp/go-build1754728840/b001/exe/main run /bin/bash
root       31202  0.0  0.0 1227308 2176 pts/1    Sl   16:20   0:00 /proc/self/exe child /bin/bash
root       31208  0.3  0.1 149668 14744 pts/1    Sl+  16:20   0:00 /usr/bin/qemu-x86_64 /bin/bash /bin/bash
root       31228  0.0  0.0   6500  3200 pts/0    R+   16:21   0:00 ps aux
```
On child process/container `ps aux` gives an error message:

```bash
Error: /proc must be mounted
  To mount /proc at boot you need an /etc/fstab line like:
      proc   /proc   proc    defaults
  In the meantime, run "mount proc /proc -t proc"
```

- The `ps` command only looks at `/proc` folder to check what all process running on the device. On the host, `/proc` is properly mounted. However, since we changed the root of the child process into a new filesystem, we have to manually mount the `/proc`.

- In the Go code, the `/proc` can be mounted as follows:

```go
func child(){

	syscall.Sethostname([]byte("container"))
	syscall.Chroot("/path/to/new/rootfs")
	syscall.Chdir("/")
	syscall.Mount("proc", "proc", "proc", 0, "") // mounting the proc

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	.
    .

	cmd.Run()

	syscall.Unmount("/proc", 0) // unmounting proc during exit
}
```
The `Mount` call accepts three arguements:
>1. **Source** : Set this to "proc" (convention). For virtual-fs, kernel usually ignores the source.
> 2. **Target** : Set this to the path to the `proc` fs. Since, we do chroot and chdir previously, simply passing `proc` would suffice.
> 3. **fsType** : must be exactly `proc`.

- Mounting `proc` would allow us to have an isolated view of process inside the container. 

- The process in the container will be still visible from host. But not the other way.

### Cloning the Mount Namespace














