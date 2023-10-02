source /etc/profile
#!/bin/bash

git pull  
go build -o visitor
pkill visitor
nohup ./visitor  &