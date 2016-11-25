package process

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mattn/go-zglob"
	"github.com/pkg/errors"
)

const (
	TuckMode   Mode = iota
	UntuckMode Mode = iota
)

type Mode int64

type Config struct {
	Source  string
	Target  string
	Modules []string
	Mode    Mode
}

func Install(config Config) error {
	// FIXME(vdemeester) use a logger system instead of fmt.Printf
	if _, err := isDirAndFollow(config.Source); err != nil {
		return err
	}
	if _, err := isDirAndFollow(config.Target); err != nil {
		return err
	}
	var visit func(base, resolvedBase, target string) filepath.WalkFunc
	switch config.Mode {
	case TuckMode:
		visit = tuckit
	case UntuckMode:
		visit = untuckit
	}
	for _, module := range config.Modules {
		matches, err := zglob.Glob(filepath.Join(config.Source, module))
		if err != nil {
			return err
		}
		for _, match := range matches {
			fmt.Println("matches", matches)
			// Kinda hacky, if the match is a symlink, follow it for match but pass it to visit
			resolvedMatch, err := isDirAndFollow(match)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Skip non folder module: %s\n", match)
				continue
			}
			fmt.Printf("tuck module: %s into %s\n", match, config.Target)
			err = filepath.Walk(resolvedMatch, visit(match, resolvedMatch, config.Target))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func tuckit(base, resolvedBase, target string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		absPath, absTargetPath, err := getAbsolutePaths(base, resolvedBase, target, path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			// Make sure the folder is user writable
			mode := info.Mode() | 0700
			err := os.Mkdir(absTargetPath, mode)
			if os.IsExist(err) {
				return nil
			}
			if err != nil {
				return err
			}
		} else {
			targetInfo, err := os.Lstat(absTargetPath)
			if err != nil && !os.IsNotExist(err) {
				return err
			}
			if targetInfo != nil {
				switch {
				case targetInfo.IsDir():
					return errors.Errorf("%s already exists and is a dir", absTargetPath)
				case targetInfo.Mode()&os.ModeSymlink != 0:
					link, err := os.Readlink(absTargetPath)
					if err != nil {
						return err
					}
					if link != absPath {
						fmt.Fprintf(os.Stderr, "%s already exists but doesn't point to %s, re-linking it", absTargetPath, absPath)
						err := os.Remove(absTargetPath)
						if err != nil {
							return err
						}
					} else {
						return nil
					}
				default:
					return errors.Errorf("%s already exists and is a plain file", absTargetPath)
				}
			}
			return os.Symlink(absPath, absTargetPath)
		}
		return nil
	}
}

func untuckit(base, resolvedBase, target string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		absPath, absTargetPath, err := getAbsolutePaths(base, resolvedBase, target, path)
		if err != nil {
			return err
		}
		targetInfo, err := os.Lstat(absTargetPath)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		if os.IsNotExist(err) || targetInfo == nil {
			fmt.Println("skipping", absTargetPath)
			return nil
		}
		if targetInfo != nil {
			switch {
			case targetInfo.IsDir():
				fmt.Println("(FIXME: should clean after) skipping", absTargetPath)
			case targetInfo.Mode()&os.ModeSymlink != 0:
				link, err := os.Readlink(absTargetPath)
				if err != nil {
					return err
				}
				if link != absPath {
					fmt.Fprintf(os.Stderr, "%s doesn't point to %s, removing it anyway", absTargetPath, absPath)
				}
				err = os.Remove(absTargetPath)
				if err != nil {
					return err
				}
			default:
				return errors.Errorf("(FIXME: should we skip) %s exists and is a plain file", absTargetPath)
			}
		}
		return nil
	}
}

func getAbsolutePaths(base, resolvedBase, target, path string) (string, string, error) {
	relativePath, err := filepath.Rel(resolvedBase, path)
	if err != nil {
		return "", "", err
	}
	targetPath := filepath.Join(target, relativePath)
	absPath, err := filepath.Abs(filepath.Join(base, relativePath))
	if err != nil {
		return "", "", err
	}
	absTargetPath, err := filepath.Abs(targetPath)
	if err != nil {
		return "", "", err
	}
	return absPath, absTargetPath, nil
}

func isDirAndFollow(dir string) (string, error) {
	fi, err := os.Lstat(dir)
	if err != nil {
		return dir, errors.Wrapf(err, "couldn't validate folder: %s", dir)
	}
	switch {
	case fi.IsDir():
		return dir, nil
	case fi.Mode()&os.ModeSymlink != 0:
		link, err := os.Readlink(dir)
		if err != nil {
			return dir, err
		}
		return isDirAndFollow(link)
	}
	return dir, errors.Errorf("%s is not a folder", dir)
}
