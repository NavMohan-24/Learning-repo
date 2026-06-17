package main

import (
	"fmt"
	"os"
	"os/exec"
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
	// cmd.SysProcAttr = &syscall.SysProcAttr{
	// 	Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	// 	Unshareflags: syscall.CLONE_NEWNS, // prevent sharing of container mount namespace to host
	// }

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	}



	must(cmd.Run()) 

}

func  child(){
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())

	syscall.Sethostname([]byte("container"))
	syscall.Chroot("/nav/rootfs")
	syscall.Chdir("/") // go to the new the root
	syscall.Mount("proc", "proc", "proc", 0, "")

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr


	must(cmd.Run()) 

	syscall.Unmount("/proc", 0)
}


func must(err error) {
    if err != nil {
        panic(err)
    }
}




