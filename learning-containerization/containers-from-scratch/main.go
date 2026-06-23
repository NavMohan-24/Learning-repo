package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

func main() {

	switch os.Args[1]{

		case "run":
			run()
		
		case "child":
			child()

		default:
			panic("bad command")
	}
}

func run(){
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())

	cmd := exec.Command("/proc/self/exe", append([]string{"child"},os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	// SysProcsAttr - struct that allow modification of OS-level process attributes
	// Clone Flags - tells the kernel what namespaces to create for new process
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		// Unshareflags: syscall.CLONE_NEWNS, // prevent sharing of container mount namespace to host
	}

	// cmd.SysProcAttr = &syscall.SysProcAttr{
	// 	Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	// }


	must(cmd.Run()) 

	//  remove cgroup upon exit
	syscall.Rmdir("/sys/fs/cgroup/nav")


}

func  child(){
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())

	cg()

	syscall.Sethostname([]byte("container"))
	syscall.Chroot("/nav/rootfs")
	syscall.Chdir("/") // go to the new the root
	syscall.Mount("proc", "proc", "proc", 0, "")

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr


	must(cmd.Run()) 
	
	// umount proc upon exit
	syscall.Unmount("/proc", 0)
}

func cg(){

	cgroups := filepath.Join("/sys/fs/cgroup/", "nav")
	
	err := os.Mkdir(cgroups, 0755)
	if err != nil && !os.IsExist(err){
		panic(err)
	}

	must(os.WriteFile(filepath.Join(cgroups, "pids.max"), []byte("20"), 0700))
	must(os.WriteFile(filepath.Join(cgroups, "cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))
}


func must(err error) {
    if err != nil {
        panic(err)
    }
}




