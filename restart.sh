git pull origin main:main
go build -o visitor
pkill visitor
nohup ./visitor  &