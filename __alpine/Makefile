# stock_ticker — Alpine (APK) packaging
#
# Recipes:
#   make build     -> abuild -r  (build the .apk containing both binaries and
#                      refresh the local repo index)
#   make client    -> dev-only: compile-check just the client binary, discarding
#                      output; skips abuild/packaging entirely
#   make daemon    -> dev-only: compile-check just the daemon binary, discarding
#                      output; skips abuild/packaging entirely
#   make release   -> upload the freshly built .apk to the Nexus apkLocal repo
#   make clean     -> run all abuild cleanup tasks
#   make info      -> print the values extracted from APKBUILD
#
# Package metadata (name / version / pkgrel / arch / binaries) is read straight
# from APKBUILD so there is a single source of truth. There is one .apk for
# both binaries: build() compiles client and daemon from the same builddir
# (one Go module, two ./client and ./daemon entrypoints), package() installs
# both to /opt/bin.

SHELL := /bin/sh

APKBUILD      := APKBUILD

NAME          := $(shell . ./$(APKBUILD); echo $$pkgname)
VERSION       := $(shell . ./$(APKBUILD); echo $$pkgver)
PKGREL        := $(shell . ./$(APKBUILD); echo $$pkgrel)
ARCH          := $(shell . ./$(APKBUILD); echo $$arch)
BINARY        := $(shell . ./$(APKBUILD); echo $$binaryname)
DAEMON_BINARY := $(shell . ./$(APKBUILD); echo $$binaryname_daemon)

NEXUSREPO := apkLocal

# abuild -r drops the package at $REPODEST/$repo/$arch/<name>-<version>-r<pkgrel>.apk,
# where $repo is the name of the directory *containing* $startdir (this project's
# checkout directory, e.g. "stock_ticker") — not pkgname. Deriving it from CURDIR
# instead of $(NAME) keeps this correct even if pkgname is renamed again.
REPO      := $(notdir $(patsubst %/,%,$(dir $(CURDIR))))
APK       := /data/packages/$(REPO)/$(ARCH)/$(NAME)-$(VERSION)-r$(PKGREL).apk

.PHONY: build client daemon release clean info

# Alpine reindex is a no-op (Nexus rebuilds APKINDEX on every upload), so the
# upload deliberately omits the -r/--reindex flag.
build:
	abuild -r

# Dev-only smoke builds; not part of the abuild/apk pipeline, and deliberately
# independent of src/*/build.sh (those are per-machine dev scripts, not meant
# to be relied on by packaging). Output is discarded — this only checks that
# the binary compiles, mirroring CGO_ENABLED=0 from APKBUILD's build().
client:
	cd ../src && CGO_ENABLED=0 go build -o /dev/null ./client

daemon:
	cd ../src && CGO_ENABLED=0 go build -o /dev/null ./daemon

release: build
	/opt/bin/nxtools upload $(NEXUSREPO) $(APK)

clean:
	abuild clean cleanpkg cleancache cleanoldpkg

info:
	@echo "NAME          = $(NAME)"
	@echo "VERSION       = $(VERSION)"
	@echo "PKGREL        = $(PKGREL)"
	@echo "ARCH          = $(ARCH)"
	@echo "BINARY        = $(BINARY)"
	@echo "DAEMON_BINARY = $(DAEMON_BINARY)"
	@echo "REPO          = $(REPO)"
	@echo "APK           = $(APK)"
