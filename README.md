battlemap-server
===

Server-side code for tabletop game automation service.

v8.go
---
v8.go requires special installation to get/build v8 code. Run `go get -d
github.com/swwu/v8.go"` to download the packages. Then cd to
`$GOPATH/src/github.com/idada/v8.go` and run `./install.sh`. This may require
installing other dependencies:

* `sudo apt-get install subversion libc6-dev g++-multilib`

May also need to force v8 to build in ia32 mode - this can be done by
modifying install.sh to run `make ia32.release` instead of `make native`, and
changing all references to `out/native` to `out/ia32.release`.


