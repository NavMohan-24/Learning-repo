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

#### Mounting a Directory

- Mounting is the process of making a storage device or partition of making storage device or partition accessible to the OS. Mounting is the process of attaching a filesystem — whether a physical storage device, network share, or virtual filesystem — to a directory (mount point), making its contents accessible at that path. This integration allows users and applications to read, write, and manage data on the mounted file system as if it were part of the local directory structure.

- Mounting a virtual fs would allow the kernel to serve the data a directory. For instance, assume we've created a directory `my_dir`, we can mount it as follows:

    `mount -t sysfs  sysfs  /my_dir   # kernel serves hardware info`

    `mount -t proc proc /my_dir # kernel serves process info`

    `mount -t devpts devpts /my_dir   # kernel serves pseudo-terminal`

- However, not all mounts serve kernel-generated data. Bind mounts simply mirror an existing path. tmpfs and overlayfs are kernel-managed but store user data rather than kernel-generated content.

#### Proc directory
- `/proc` is a virtual filesystem used by the system during runtime. It holds information of all the process that are currently running on the system. 

- The `ps` command looks at `/proc` folder to check what all process running on the device. 

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
<!-- - Ideally containers should be isolated from the host and should only show the process within it. -->

- On the host, `ps` command returns all the process since the `/proc` is properly mounted. It includes process inside the container too.

- However, inside the container, it returns an error message. Ideally, it should provide an isolated view of all the process that are inside the container. 

- It happens since we haven't mounted the `/proc` directory in our new filesystem of the container (to which we `chroot`-ed). 


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

- Mounting `proc` would allow us to have an isolated view of process inside the container. The process in the container will be still visible from host. But not the other way.



### Cloning the Mount Namespace

- Running `mount` command on the host and container gives the result as follows:

ON HOST:

```bash
root@0ff0962a3623:/#  mount | grep proc
proc on /proc type proc (rw,nosuid,nodev,noexec,relatime)
proc on /nav/rootfs/proc type proc (rw,relatime)
```
ON CONTAINER:

```bash
mount | grep proc
proc on /proc type proc (rw,relatime)
```
- Here we see the container mount polluting the host's mount table. i.e, mount point of the containers proc is visible on the host.

    > A mount table is an internal list maintained by the Linux kernel that keeps track of every active filesystem on the system. It maps a filesystem source (like a hard drive partition, a network share, or a virtual driver like procfs) to its corresponding destination directory (the mount point).

- It occurs since the container we have created shares the same `mount` namespace. This poses a security threat as the mount points of host could be modified from the container.

- Ideally, we wish the container to have isolated mount points from the host. This could be done by setting a new mount namespace for a container. 

- The container would maintain a seperate mount table when the mount namespace is set.

- In the GO code a new mount namespace could be created as follows:

``` go
func run(){
    cmd := exec.Command("/proc/self/exe", append([]string{"child"},os.Args[2:]...)...)
    .
    .

    cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		// Unshareflags: syscall.CLONE_NEWNS, 
	}

    cmd.Run()
    }
```


- When a new mount namespace is created with `CLONE_NEWNS`, the kernel copies the host's mount table into the container's mount namespace as a starting point — including the host's `/proc`. However, after `chroot` shifts the root to `/nav/rootfs`, the `child()` function remounts `/proc` fresh inside the new root. This new `/proc` is scoped to the container's PID namespace, so it only reflects the container's processes — overriding the copied one.






| Path | Filesystem Type | Purpose |
|------|----------------|---------|
| `/proc` | `procfs` | Process information |
| `/sys` | `sysfs` | Hardware and kernel information |
| `/dev` | `devtmpfs` | Device files |
| `/dev/pts` | `devpts` | Pseudo-terminals |











