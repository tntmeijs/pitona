# Ensure a clean build folder is available
echo "Creating build directory"
rm -rf build
mkdir build

# Build the server executable for
export GOOS="linux"
export GOARCH=arm
export GOARM=7

pushd server
echo "Building server application"
go build -buildmode=exe -o ../build/pitona main.go
popd

# Build the website
echo "Building client application"
pushd client
npm run build
popd

# Move the website files from the client to the server
echo "Moving client application data to build folder"
mv client/build build/public

# Pack everything up into an archive to make it easier to move around
pushd build
tar -czvf pitona.tar.gz *
rm -rf public
rm pitona
popd
