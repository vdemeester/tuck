# 🐄 tuck [![Build Status](https://travis-ci.org/vdemeester/tuck.svg?branch=master)](https://travis-ci.org/vdemeester/tuck)

`tuck` is symlink farm manager *à-la-stow* written in Go. In a
nutshell, it's a tools that create symlink in a target based on a
source folder and modules.

```bash
$ tree testdata
testdata
├── a
│   └── b
│       └── 1
├── b
│   └── 2
├── c
│   ├── 5
│   └── d
│       └── 3
└── d
    ├── 6
    ├── e
    │   ├── 7
    │   └── f
    │       └── 4
    └── f
        ├── 9
        └── g
            └── 8

10 directories, 9 files
$ tree target
target

0 directories, 0 files

# the following copies the content of the folder a from testdata to target
$ tuck -s testdata -t target a
tuck module: testdata/a into target
$ tree target
target
└── b
    └── 1 -> /home/vincent/go/src/github.com/vdemeester/tuck/testdata/a/b/1

1 directory, 1 file

# it supports globing
$ tuck -s testdata -t target '*'
tuck module: testdata/a into target
tuck module: testdata/b into target
tuck module: testdata/c into target
tuck module: testdata/d into target
$ tree target
target
├── 2 -> /home/vincent/go/src/github.com/vdemeester/tuck/testdata/b/2
├── 5 -> /home/vincent/go/src/github.com/vdemeester/tuck/testdata/c/5
├── 6 -> /home/vincent/go/src/github.com/vdemeester/tuck/testdata/d/6
├── b
│   └── 1 -> /home/vincent/go/src/github.com/vdemeester/tuck/testdata/a/b/1
├── d
│   └── 3 -> /home/vincent/go/src/github.com/vdemeester/tuck/testdata/c/d/3
├── e
│   ├── 7 -> /home/vincent/go/src/github.com/vdemeester/tuck/testdata/d/e/7
│   └── f
│       └── 4 -> /home/vincent/go/src/github.com/vdemeester/tuck/testdata/d/e/f/4
└── f
    ├── 9 -> /home/vincent/go/src/github.com/vdemeester/tuck/testdata/d/f/9
    └── g
        └── 8 -> /home/vincent/go/src/github.com/vdemeester/tuck/testdata/d/f/g/8

6 directories, 9 files
```

It's also possible to remove the symlink the same way, just use `-D`
or `--delete` flag. There is still a `TODO` on cleaning empty folders
afterwards though.

```bash
$ tuck -d testdata -t target -D '*'                                                                                                                            ~/go/src/github.com/vdemeester/tuck
tuck module: testdata/a into target
(FIXME: should clean after) skipping /home/vincent/go/src/github.com/vdemeester/tuck/target
(FIXME: should clean after) skipping /home/vincent/go/src/github.com/vdemeester/tuck/target/b
tuck module: testdata/b into target
(FIXME: should clean after) skipping /home/vincent/go/src/github.com/vdemeester/tuck/target
tuck module: testdata/c into target
(FIXME: should clean after) skipping /home/vincent/go/src/github.com/vdemeester/tuck/target
(FIXME: should clean after) skipping /home/vincent/go/src/github.com/vdemeester/tuck/target/d
tuck module: testdata/d into target
(FIXME: should clean after) skipping /home/vincent/go/src/github.com/vdemeester/tuck/target
(FIXME: should clean after) skipping /home/vincent/go/src/github.com/vdemeester/tuck/target/e
(FIXME: should clean after) skipping /home/vincent/go/src/github.com/vdemeester/tuck/target/e/f
(FIXME: should clean after) skipping /home/vincent/go/src/github.com/vdemeester/tuck/target/f
(FIXME: should clean after) skipping /home/vincent/go/src/github.com/vdemeester/tuck/target/f/g
$ tree target
target
├── b
├── d
├── e
│   └── f
└── f
    └── g

6 directories, 0 files
```
