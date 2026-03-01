#
# NOTE:
#
# The Makefile is legacy and should not be used for building;
# instead, use the installer delegated for your OS in the build directory.
#
# The Makefile should only be used for building legacy Luna L1 software OR to build the MacOS installer
#
SRC=./

all: bin/luna-l2 bin/las bin/l2ld bin/lcc bin/lcc1
.PHONY: clean install

bin/luna-l2: $(SRC)/l2/*
	cd l2 && go build -o ../bin/luna-l2 ./luna_l2.go	

bin/las: $(SRC)/las/*
	cd las && go build -o ../bin/las ./las.go

bin/lcc1: $(SRC)/lcc1/*
	cd lcc1 && go build -o ../bin/lcc1 ./lcc1.go

bin/lcc: $(SRC)/lcc/*
	cd lcc && go build -o ../bin/lcc ./lcc.go

bin/l2ld: $(SRC)/l2ld/*	
	cd l2ld && go build -o ../bin/l2ld ./l2ld.go

macos-installer:
	cd l2 && CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o ../Mac/amd64/usr/local/bin/"Luna L2"/Contents/MacOS/luna-l2 luna_l2.go
	cd lcc && GOOS=darwin GOARCH=amd64 go build -o ../Mac/amd64/usr/local/bin/lcc lcc.go
	cd las && GOOS=darwin GOARCH=amd64 go build -o ../Mac/amd64/usr/local/bin/las las.go
	cd lcc1 && GOOS=darwin GOARCH=amd64 go build -o ../Mac/amd64/usr/local/bin/lcc1 lcc1.go
	cd l2ld && GOOS=darwin GOARCH=amd64 go build -o ../Mac/amd64/usr/local/bin/l2ld l2ld.go
	cd l2 && CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o ../Mac/arm64/usr/local/bin/"Luna L2"/Contents/MacOS/luna-l2 luna_l2.go
	cd lcc && GOOS=darwin GOARCH=arm64 go build -o ../Mac/arm64/usr/local/bin/lcc lcc.go
	cd las && GOOS=darwin GOARCH=arm64 go build -o ../Mac/arm64/usr/local/bin/las las.go
	cd lcc1 && GOOS=darwin GOARCH=arm64 go build -o ../Mac/arm64/usr/local/bin/lcc1 lcc1.go
	cd l2ld && GOOS=darwin GOARCH=arm64 go build -o ../Mac/arm64/usr/local/bin/l2ld l2ld.go
	pkgbuild \
		--root Mac/amd64 \
		--install-location / \
		--identifier com.alexfdev0.lunal2.amd64 \
		--version 1.0 \
		--scripts Mac/scripts \
		build/"Luna L2 (amd64).pkg"
	pkgbuild \
		--root Mac/arm64 \
		--install-location / \
		--identifier com.alexfdev0.lunal2.arm64 \
		--version 1.0 \
		--scripts Mac/scripts \
		build/"Luna L2 (arm64).pkg"

windows-installer:
	cd l2 && CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o ../Windows/luna-l2.exe luna_l2.go
	cd lcc && GOOS=windows GOARCH=amd64 go build -o ../Windows/lcc.exe lcc.go
	cd las && GOOS=windows GOARCH=amd64 go build -o ../Windows/las.exe las.go
	cd lcc1 && GOOS=windows GOARCH=amd64 go build -o ../Windows/lcc1.exe lcc1.go
	cd l2ld && GOOS=windows GOARCH=amd64 go build -o ../Windows/l2ld.exe l2ld.go
	cd Windows && wixl -v msi.xml -o "../build/Luna L2.msi"

install:
	sudo cp bin/* /usr/local/bin/

clean:
	rm -f /usr/local/bin/luna-l2
	rm -f /usr/local/bin/las
	rm -f /usr/local/bin/lcc1
	rm -f /usr/local/bin/lcc
	rm -f /usr/local/bin/l2ld
