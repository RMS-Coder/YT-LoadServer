package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

var fileModTimes = make(map[string]time.Time)
var currentProc *exec.Cmd

func runApp() {
	fmt.Println("\nRunning...")

	if currentProc != nil {
		err := currentProc.Process.Kill()

		if err != nil {
			fmt.Printf("Error terminating process: %v\n", err)
		}
		currentProc.Wait()
	}

	currentProc = exec.Command("./app.exe")

	currentProc.Stdout = os.Stdout
	currentProc.Stderr = os.Stderr

	err := currentProc.Start()
	if err != nil {
		fmt.Printf("Error starting app: %v\n", err)
		currentProc = nil
	}

	fmt.Println("\nReady!")
}

func rebuildApp() {
	fmt.Println("\n=== Recompiling APP ===")

	cmd := exec.Command("go", "build", "-o", "../app.exe")
	cmd.Dir = "./app"

	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("Error generating APP: %v\n", err)
		fmt.Printf("Compiler output:\n%s\n", output)
	}
}

func checkChanges() bool {
	changed := false

	err := filepath.Walk("app", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		lastMod, exists := fileModTimes[path]
		if !exists || info.ModTime().After(lastMod) {
			fileModTimes[path] = info.ModTime()
			changed = true
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error reading directory:", err)
	}

	return changed
}

func main() {
	fmt.Println("Go Live Reload - Starting monitoring...")
	fmt.Println("Monitoring folder: app/")
	fmt.Println("Press Ctrl C to exit")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-interrupt
		fmt.Println("\nTurning off...")
		if currentProc != nil {
			currentProc.Process.Kill()
		}
		os.Exit(0)
	}()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if checkChanges() {
			fmt.Println("\n---------------------------------")
			rebuildApp()
			runApp()
		}
	}
}
