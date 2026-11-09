App := 'makit'
Version := `sed -n 's/^const VERSION = "\(.*\)"/\1/p' cmd/main/version.go`

# show help message
@default: help

# show help message
@help:
    echo "Build tool for {{ App }} {{ Version }} with Just"
    echo "Usage: just <recipe>"
    echo ""
    just --list

# build the application with running tests
build: test
    go build -o makit ./cmd/main

# run tests and generate the coverage report
test:
    go test -covermode=count -coverprofile=coverage.out ./...

# clean up build artifacts
clean:
    go clean
    rm -rf coverage.out makit makit.exe dist

# update the version if the new version is provided
update_version new_version = "":
    if [ "{{ new_version }}" != "" ]; then \
        sed 's/$VERSION/{{ new_version }}/g' .template/README.md > README.md; \
        sed 's/$VERSION/{{ new_version }}/g' .template/version.go > cmd/main/version.go; \
    fi

# build makit for all platforms
make_distribution_files:
    for os in "linux" "windows" "darwin"; do \
        for arch in "amd64" "arm64"; do \
            bin=makit; \
            if [ "$os" = "windows" ]; then bin=makit.exe; fi; \
            mkdir -p dist/makit-$os-$arch; \
            env GOOS=$os GOARCH=$arch CGO_ENABLED=0 go build -o dist/makit-$os-$arch/$bin ./cmd/main; \
            cp README.md LICENSE dist/makit-$os-$arch; \
            tar cvfz dist/makit-$os-$arch.tar.gz -C dist makit-$os-$arch; \
        done; \
    done

# upload distribution files to the GitHub release
upload_assets tag:
    gh release upload --repo k140918/makit {{ tag }} dist/*.tar.gz
