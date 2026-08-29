set shell := ["sh", "-eu", "-c"]

configure:
    cmake -S . -B build/release -G Ninja -DCMAKE_BUILD_TYPE=Release

build: configure
    cmake --build build/release
    cmake --install build/release --prefix dist

run: build
    ./dist/bin/hanji

clean:
    cmake -E remove_directory build
    cmake -E remove_directory dist
