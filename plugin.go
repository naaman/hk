package main

import (
	"log"
	"net/url"
	"os"
	"os/exec"
	"syscall"
)

var hkPath string

func init() {
	const defaultPluginPath = "/usr/local/lib/hk/plugin"
	hkPath = os.Getenv("HKPATH")
	if hkPath == "" {
		hkPath = defaultPluginPath
	}

}

func execPlugin(path string, args []string) error {
	u, err := url.Parse(apiURL)
	if err != nil {
		log.Fatal(err)
	}

	hkuser, hkpass := getCreds(u)
	u.User = url.UserPassword(hkuser, hkpass)
	hkapp, _ := app()
	env := []string{
		"HEROKU_API_URL=" + u.String(),
		"HKAPP=" + hkapp,
		"HKUSER=" + hkuser,
		"HKPASS=" + hkpass,
		"HKHOST=" + u.Host,
		"HKVERSION=" + Version,
	}

	return syscall.Exec(path, args, append(env, os.Environ()...))
}

func findPlugin(name string) (path string) {
	path = lookupPlugin(name)
	if path == "" {
		path = lookupPlugin("default")
	}
	return path
}

// NOTE: lookupPlugin is not threadsafe for anything needing the PATH env var.
func lookupPlugin(name string) string {
	opath := os.Getenv("PATH")
	defer os.Setenv("PATH", opath)
	os.Setenv("PATH", hkPath)

	path, err := exec.LookPath(name)
	if err != nil {
		if e, ok := err.(*exec.Error); ok && e.Err == exec.ErrNotFound {
			return ""
		}
		log.Fatal(err)
	}
	return path
}
