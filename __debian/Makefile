# stock_ticker — Debian (.deb) packaging
#
# Recipes:
#   make build     -> build both binaries and assemble the .deb
#   make client    -> dev-only: compile-check just the client binary, discarding
#                      output; skips the dpkg-deb pipeline entirely
#   make daemon    -> dev-only: compile-check just the daemon binary, discarding
#                      output; skips the dpkg-deb pipeline entirely
#   make release   -> upload the freshly built .deb to the Nexus aptLocal repo (with reindex)
#   make clean     -> remove the staging tree and the built .deb
#   make info      -> print the values extracted from the control file
#
# Package metadata (name / version / arch) is read straight from the debian
# control file so there is a single source of truth. Binary names are not:
# the control file is Debian-policy metadata, not a script, so there's nowhere
# to put a "binaryname" var the way APKBUILD has one — BINARY/DAEMON_BINARY
# are declared here instead. This build recipe folds together what
# 2.build_binary.sh + restore_repo.sh used to do: it copies (never moves) the
# DEBIAN control files into a throwaway staging tree, so there is nothing to
# git-restore afterwards.

SHELL       := /bin/bash
.SHELLFLAGS := -o pipefail -c

CONTROL   := control

NAME          := $(shell awk -F': +' '/^Package:/{print $$2}'      $(CONTROL))
VERSION       := $(shell awk -F': +' '/^Version:/{print $$2}'      $(CONTROL))
ARCH          := $(shell awk -F': +' '/^Architecture:/{print $$2}' $(CONTROL))
BINARY        := stockticker
DAEMON_BINARY := $(BINARY)Daemon

NEXUSREPO := aptLocal
SRCDIR    := ../src
STAGING   := $(NAME)-$(VERSION)_$(ARCH)
DEB       := $(CURDIR)/$(STAGING).deb

.PHONY: build client daemon release upload clean info

build:
	rm -rf "$(STAGING)"
	install -d "$(STAGING)/DEBIAN" "$(STAGING)/opt/bin" "$(STAGING)/opt/sbin"
	cp control preinst prerm postinst postrm "$(STAGING)/DEBIAN/"
	cd $(SRCDIR) && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -buildid=" \
	  -o "$(CURDIR)/$(STAGING)/opt/bin/$(BINARY)" ./client
	cd $(SRCDIR) && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -buildid=" \
	  -o "$(CURDIR)/$(STAGING)/opt/sbin/$(DAEMON_BINARY)" ./daemon
	dpkg-deb -b "$(STAGING)" "$(DEB)"

# Dev-only smoke builds; not part of the dpkg-deb pipeline, and deliberately
# independent of src/*/build.sh (those are per-machine dev scripts, not meant
# to be relied on by packaging). Output is discarded.
client:
	cd $(SRCDIR) && CGO_ENABLED=0 go build -o /dev/null ./client

daemon:
	cd $(SRCDIR) && CGO_ENABLED=0 go build -o /dev/null ./daemon

upload:
	/opt/bin/nxtools upload -r $(NEXUSREPO) $(DEB)


clean:
	rm -rf "$(STAGING)" "$(DEB)"

release:
	${MAKE} build
	${MAKE} upload
	${MAKE} clean

info:
	@echo "NAME          = $(NAME)"
	@echo "VERSION       = $(VERSION)"
	@echo "ARCH          = $(ARCH)"
	@echo "BINARY        = $(BINARY)"
	@echo "DAEMON_BINARY = $(DAEMON_BINARY)"
	@echo "STAGING       = $(STAGING)"
	@echo "DEB           = $(DEB)"
