# Usage of simple-ssh

####Build
    go build -o simple-ssh cmd/simple-ssh/main.go
    
####Example of usage    	
    ./simple-ssh -h 192.168.1.10 -u user
    
####To view the arguments
    ./simple-ssh --help
    
####Table of arguments
| Argument       | Type | Assignment of the argument | Default |
| ------------- | ------------- | ------------- | -------------|
| -h  | String                  | Host          |              |
| -p  | Int                     | Port          | 22           |
| -pk | String                  | Private key   | $HOME/.ssh/id_rsa |
| -u | String                   | User name     |             |