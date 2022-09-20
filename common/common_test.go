package common

import (
	"fmt"
	"testing"
)

func TestInit(t *testing.T) {

}

func TestUtils(t *testing.T) {
	
	fmt.Printf("mydocker path: %v\n",GetMyDockerPath())
	fmt.Printf("proc path: %v\n",GetProcPath())
	fmt.Printf("cgroup path: %v\n",GetCgroupPath())

}