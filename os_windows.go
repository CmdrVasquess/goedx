package goedx

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

var relJournalPath = []string{
	"",
	"Saved Games",
	"Frontier Developments",
	"Elite Dangerous",
}

func FindJournals() (dir string, err error) {
	usr, err := user.Current()
	if err != nil {
		return ".", err
	}
	relJournalPath[0] = usr.HomeDir
	res := filepath.Join(relJournalPath...)
	if stat, err := os.Stat(res); err != nil {
		return "", err
	} else if !stat.IsDir() {
		return "", fmt.Errorf("'%s' is not a directory", res)
	}
	return res, nil
}
