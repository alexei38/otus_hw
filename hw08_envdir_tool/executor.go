package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func env2map(envs []string) map[string]string {
	m := make(map[string]string)
	for _, e := range envs {
		if i := strings.Index(e, "="); i >= 0 {
			m[e[:i]] = e[i+1:]
		}
	}
	return m
}

func map2env(envs map[string]string) []string {
	m := []string{}
	for k, v := range envs {
		m = append(m, strings.Join([]string{k, v}, "="))
	}
	return m
}

func mergeEnv(env Environment) []string {
	mEnv := env2map(os.Environ())
	for k, v := range env {
		_, ok := mEnv[k]
		if ok {
			delete(mEnv, k)
		}
		if !v.NeedRemove {
			mEnv[k] = v.Value
		}
	}
	return map2env(mEnv)
}

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		fmt.Fprintln(os.Stderr, "not enough arguments")
		return 111
	}
	c := exec.Command("/usr/bin/env", "bash", "-c", strings.Join(cmd, " ")) // nolint:gosec
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	c.Env = mergeEnv(env)

	if err := c.Run(); err != nil {
		exitErr, ok := err.(*exec.ExitError) // nolint:errorlint
		if ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus()
			}
		}
		// Если неизвестная ошибка, напишем в stderr и вернем 111 код
		fmt.Fprintln(os.Stderr, err)
		return 111
	}
	return 0
}
