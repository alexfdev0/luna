mkdir -p ../bin
cd ../l2 && go build -o ../bin/luna-l2 ./luna_l2.go	
cd ../las && go build -o ../bin/las ./las.go
cd ../lcc1 && go build -o ../bin/lcc1 ./lcc1.go
cd ../lcc && go build -o ../bin/lcc ./lcc.go
cd ../l2ld && go build -o ../bin/l2ld ./l2ld.go
cp ../bin/* /usr/local/bin
mkdir -p /usr/local/lib/lcc
mkdir -p /usr/local/lib/l2ld
cp -n ../l2ld/libs.conf /usr/local/lib/l2ld/
