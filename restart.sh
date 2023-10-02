source /etc/profile
#!/bin/bash

sudo git pull  
sudo go build -o visitor
sudo pkill visitor
nohup ./visitor  &