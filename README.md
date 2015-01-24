battlemap-server
===

Server-side code for tabletop game automation service.

Installing
---

### v8.go
v8.go requires special installation to get/build v8 code. Run `go get -d
github.com/swwu/v8.go"` to download the packages. Then cd to
`$GOPATH/src/github.com/idada/v8.go` and run `./install.sh`. This may require
installing other dependencies:

* `sudo apt-get install subversion libc6-dev g++-multilib`

May also need to force v8 to build in ia32 mode - this can be done by
modifying install.sh to run `make ia32.release` instead of `make native`, and
changing all references to `out/native` to `out/ia32.release`.


Object Model
---

The basic character-sheet evaluation consists of five primitives. These are
entities, effects, rules, variables, and reductions.

### Entity
An entity is a data structure which represents the character sheet and its
underlying character. An entity stores, as its state, all active effects, and
some data variables representing base stats.

The entity will evaluate the union of the rules attached to all its effects
when asked to generate its instantaneous state.

### Effect
An effect represents a grouping of rules that form a coherent mechanic.

### Rule
Rules represent a grouping of variable declarations and reductions.

### Reduction
A reduction represents a complete bipartite dependency subgraph between
variables, which modifies the value of some set of variables based on the
values of some set of dependencies.

Reductions are evaluated in dependency ordering, meaning that a reduction
which depends on a variable will always be evaluated after all reductions
which modify that variable.

### Variable
A variable represents a single value. There are two types of variables - data
variables and reducer variables.

A data variable is simply a fixed value which is provided via user input (e.g.
base stats).

A reducer variable is a variable which contains a well-defined set of
operations which can be called, and which reduces those operations in a
deterministic order regardless of the order they are called in.

