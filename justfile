set shell := ["sh", "-eu", "-c"]

configure:
    cmake -S . -B build/release -G Ninja -DCMAKE_BUILD_TYPE=Release

build: configure
    cmake --build build/release
    cmake --install build/release --prefix dist

# Build an unsigned Debian binary package under build/debian.
package:
    mkdir -p build/debian; \
    package_version="$(dpkg-parsechangelog -ldebian/changelog --show-field Version)"; \
    package_arch="$(dpkg-architecture --query DEB_HOST_ARCH)"; \
    package_basename="hanji_${package_version}_${package_arch}"; \
    trap 'debian/rules clean >/dev/null' EXIT; \
    dpkg-buildpackage --build=binary --no-sign \
        --buildinfo-file="build/debian/${package_basename}.buildinfo" \
        --buildinfo-option=-ubuild/debian \
        --changes-file="build/debian/${package_basename}.changes" \
        --changes-option=-ubuild/debian

run: build
    ./dist/bin/hanji

clean:
    cmake -E remove_directory build
    cmake -E remove_directory dist
