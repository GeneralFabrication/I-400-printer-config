package main

import (
    "fmt"
    "io"
    "os"
    "os/exec"
)

func runCommand(name string, args ...string) error {
    cmd := exec.Command(name, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    if err != nil {
        return fmt.Errorf("error running command %s %v: %v", name, args, err)
    }
    return nil
}

func copyFile(src, dest string) error {
    from, err := os.Open(src)
    if err != nil {
        return fmt.Errorf("error opening source file %s: %w", src, err)
    }
    defer from.Close()

    to, err := os.Create(dest)
    if err != nil {
        return fmt.Errorf("error creating destination file %s: %w", dest, err)
    }
    defer to.Close()

    _, err = io.Copy(to, from)
    if err != nil {
        return fmt.Errorf("error copying from %s to %s: %w", src, dest, err)
    }

    return nil
}
