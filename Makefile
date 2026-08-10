.PHONY: run build package

DEBIAN_OUTPUT := build/debian
PACKAGE_VERSION := $(shell dpkg-parsechangelog --show-field Version)
PACKAGE_ARCH := $(shell dpkg-architecture --query DEB_HOST_ARCH)
PACKAGE_BASENAME := hanji_$(PACKAGE_VERSION)_$(PACKAGE_ARCH)

run:
	go run .

build:
	go build -o ./build/hanji .

# Build an unsigned Debian binary package under build/debian.
package:
	@mkdir -p $(DEBIAN_OUTPUT)
	@trap '$(MAKE) -f debian/rules clean >/dev/null' EXIT; \
		dpkg-buildpackage --build=binary --no-sign \
			--buildinfo-file=$(DEBIAN_OUTPUT)/$(PACKAGE_BASENAME).buildinfo \
			--buildinfo-option=-u$(DEBIAN_OUTPUT) \
			--changes-file=$(DEBIAN_OUTPUT)/$(PACKAGE_BASENAME).changes \
			--changes-option=-u$(DEBIAN_OUTPUT)
