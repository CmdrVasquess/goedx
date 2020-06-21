package edgx

import (
	"fmt"
	"os"
)

var relJournalPath = []string{
	"",
	"Saved Games",
	"Frontier Developments",
	"Elite Dangerous",
}

func FindJournals() (dir string, err error) {
	if usr, err := user.Current(); err != nil {
		return "."
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
