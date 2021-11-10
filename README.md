# git-get

Clone git repository in $GIT_GET_PROJECTS_PATH and preserves `<githost>/<username>/<repository>` structure.

## Installation

### Homebrew

```bash
brew install pixelfactoryio/tools/git-get
```

### Usage

```bash
$ export GIT_GET_PROJECTS_PATH=$HOME/Projects
$ git get https://github.com/pixelfactoryio/git-get.git
Cloning into '/Users/amine/Projects/github.com/pixelfactoryio/git-get'...
remote: Enumerating objects: 83, done.
remote: Counting objects: 100% (83/83), done.
remote: Compressing objects: 100% (44/44), done.
remote: Total 83 (delta 26), reused 78 (delta 26), pack-reused 0
Receiving objects: 100% (83/83), 6.60 MiB | 1.50 MiB/s, done.
Resolving deltas: 100% (26/26), done.
```

### Directory structure

git-get will parse the repository URL and will create the following directory structure using `$GIT_GET_PROJECTS_PATH` as root directory.

```bash
.
└── github.com
    └── <username>
        └── <repository>
```
