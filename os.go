// +build !windows

package edgx

func FindJournals() (dir string, err error) { return ".", nil }
