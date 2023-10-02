source /etc/profile
#!/bin/bash

git pull  
go build main.go -o visitor
pkill visitor
nohup ./visitor  &