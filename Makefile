#
# NOTE:
#
# The Makefile is legacy and should not be used for building;
# instead, use the installer delegated for your OS in the build directory.
#
# The Makefile should only be used for building legacy Luna L1 software OR to build the MacOS installer
#

LUAC=env luac
LUA=env lua
SRC=./

all: bin/luna-l2 bin/las bin/l2ld bin/lcc bin/lcc1 
legacy: luna-l1
.PHONY: clean install

luna-l1:
	sudo mkdir -p /usr/local/bin/lvm
	sudo $(LUAC) -o /usr/local/bin/lvm/luna-l1 $(SRC)/l1/luna_l1.lua
	sudo printf '#!/bin/sh\n $(LUA) /usr/local/bin/lvm/luna-l1 "$$@"' > /usr/local/bin/luna-l1
	sudo $(LUAC) -o /usr/local/bin/lvm/lcc-l1 $(SRC)/l1/lcc.lua
	sudo printf '#!/bin/sh\n $(LUA) /usr/local/bin/lvm/lcc-l1 "$$@"' > /usr/local/bin/lcc-l1
	sudo $(LUAC) -o /usr/local/bin/lvm/lasm-l1 $(SRC)/l1/lasm.lua
	sudo printf '#!/bin/sh\n $(LUA) /usr/local/bin/lvm/lasm-l1 "$$@"' > /usr/local/bin/lasm-l1	
	sudo chmod +x /usr/local/bin/luna-l1
	sudo chmod +x /usr/local/bin/lcc-l1
	sudo chmod +x /usr/local/bin/lasm-l1

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
	sudo cp bin/luna-l2 Mac/pkgroot2/usr/local/bin/"Luna L2"/Contents/MacOS/
	sudo cp bin/lcc Mac/pkgroot2/usr/local/bin/
	sudo cp bin/las Mac/pkgroot2/usr/local/bin/
	sudo cp bin/lcc1 Mac/pkgroot2/usr/local/bin/
	sudo cp bin/l2ld Mac/pkgroot2/usr/local/bin/	
	pkgbuild \
		--root Mac/pkgroot2 \
		--install-location / \
		--identifier com.alexfdev0.lunal2.tools \
		--version 1.0 \
		--scripts Mac/scripts \
		build/"Luna L2.pkg"

windows-installer:
	cd l2 && GOOS=windows GOARCH=amd64 go build -o ../Windows/luna-l2.exe luna_l2.go
	cd lcc && GOOS=windows GOARCH=amd64 go build -o ../Windows/lcc.exe lcc.go
	cd las && GOOS=windows GOARCH=amd64 go build -o ../Windows/las.exe las.go
	cd lcc1 && GOOS=windows GOARCH=amd64 go build -o ../Windows/lcc1.exe lcc1.go
	cd l2ld && GOOS=windows GOARCH=amd64 go build -o ../Windows/l2ld.exe l2ld.go
	cd Windows && wixl -v msi.xml -o "Luna L2.msi"	

clean:
	rm -rf /usr/local/bin/lvm
	rm -f /usr/local/bin/luna-l1
	rm -f /usr/local/bin/lasm-l1
	rm -f /usr/local/bin/lcc-l1
	rm -f /usr/local/bin/luna-l2
	rm -f /usr/local/bin/las
	rm -f /usr/local/bin/lcc1
	rm -f /usr/local/bin/lcc
	rm -f /usr/local/bin/l2ld
