mkdir bin
cd l2 && go build -o ..\bin\luna-l2.exe luna_l2.go && cd ..\
cd l2ld && go build -o ..\bin\l2ld.exe l2ld.go && cd ..\
cd las && go build -o ..\bin\las.exe las.go && cd ..\
cd lcc && go build -o ..\bin\lcc.exe lcc.go && cd ..\
cd lcc1 && go build -o ..\bin\lcc1.exe lcc1.go && cd ..\
mkdir C:\"Program Files (x86)"\"Luna L2"\
mkdir C:\"Program Files (x86)"\"Luna L2"\bin
mkdir C:\"Program Files (x86)"\"Luna L2"\lib
mkdir C:\"Program Files (x86)"\"Luna L2"\lib\lcc
mkdir C:\"Program Files (x86)"\"Luna L2"\lib\l2ld
xcopy bin\* C:\"Program Files (x86)"\"Luna L2"\bin\ /Y /I
type nul > C:\"Program Files (x86)"\"Luna L2"\lib\l2ld\libs.conf
setx PATH "%PATH%;C:\Program Files (x86)\Luna L2\bin" /M
