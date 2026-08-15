//go:build !windows

package main

import "fmt"

func showConfigWindow(_ string, _ appConfig) (bool, error) {
	return false, fmt.Errorf("the graphical setup is currently available only on Windows")
}

func showFatalError(message string) {
	fmt.Println("EasyTracker:", message)
}
