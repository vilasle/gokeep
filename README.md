#### Gokeep

Client-server application for keeping private information.

Server use PostgreSQL how storage 

#### Running server
```bash
git clone https://github.com/vilasle/gokeep.git
cd gokeep
go mod tidy
make server
#create certificates
./rsa.sh

#create salt for 
sha1pass bin/gokeep-backend > cert/salt   
./bin/gokeep-backend \
	--addr :9092
	--db-url 'postgres://postgres:password@host:5432/db?sslmode=disable' \
	--salt ./cert/salt \
	--certificate ./cert/server.crt \
	--private-key ./cert/server.key
```


#### Build client
```bash
git clone https://github.com/vilasle/gokeep.git
cd gokeep
make client 
```
#### Init configuration
```bash
./client --workspace /path/to/workspace \
	--ca-file cert/ca.crt \
	config init --gprc-socket :9092 --db-path /path/to/local/db
```
#### Registration and login
```bash
#create account
./client --workspace /path/to/workspace \
	--ca-file cert/ca.crt \
	account create $account_name

#login
./client --workspace /path/to/workspace \
	--ca-file cert/ca.crt \
	account create $account_name

```
##### Work with private data


##### Login and passwords
```bash
#add
./client --workspace /path/to/workspace \
	--ca-file cert/ca.crt 
	data cred add --login login --password 1234567890 --metadata scope=test
#update  
./client --workspace /path/to/workspace \
	--ca-file cert/ca.crt 
	data cred edit --id 1 --login login --password 1234567890
#get all 
./client --workspace /path/to/workspace \
	--ca-file cert/ca.crt 
	data cred get
#get specific
./client --workspace /path/to/workspace \
	--ca-file cert/ca.crt 
	data cred get --id 1
#delete specific 
./client --workspace /path/to/workspace \
	--ca-file cert/ca.crt 
	data cred delete --id 1
```
##### Bank cards
```bash
#add
./client --workspace /path/to/workspace \
	--ca-file cert/ca.crt 
	data bank add --number 1234567890 --cvv 432 --expires '03/31' --metadata scope=test
#update  
./client --workspace /path/to/workspace \
	--ca-file cert/ca.crt 
	data bank edit --id 1 --number 1234567890 --cvv 432 --expires '03/31' --metadata scope=test 
#get all 
./client --workspace /path/to/workspace \
	--ca-file cert/ca.crt 
	data bank get
#get specific cred
./client --workspace /path/to/workspace \
	--ca-file cert/ca.crt 
	data bank get --id 1
#delete specific cred 
./client --workspace /path/to/workspace \
	--ca-file cert/ca.crt 
	data bank delete --id 1
```

##### Text data
```bash
#add
./client --workspace /path/to/workspace \
	--ca-file cert/ca.crt 
	data text add --path /path/to/file --name some_note --metadata scope=test
	# or
	data text add --data "some test data" --name some_note --metadata scope=test
#update  
./client --workspace /path/to/workspace \
	--ca-file cert/ca.crt 
	text edit --id 1 --path /path/to/file --name some_note --metadata scope=test
	# or
	data text edit --id 1 --data "another text" --name some_note --metadata scope=test
#get all 
./client --workspace /path/to/workspace \
	--ca-file cert/ca.crt 
	data text get
#get specific cred
./client --workspace /path/to/workspace \
	--ca-file cert/ca.crt 
	data text get --id 1
#delete specific cred 
./client --workspace /path/to/workspace \
	--ca-file cert/ca.crt 
	data text delete --id 1
```
##### Binary data
```bash
#add
./client --workspace /path/to/workspace \
	--ca-file cert/ca.crt 
	data binary add --path /path/to/file --name some_note --metadata scope=test
#update  
./client --workspace /path/to/workspace \
	--ca-file cert/ca.crt 
	data text binary --id 1 --path /path/to/file --name some_note --metadata scope=test
#get all 
./client --workspace /path/to/workspace \
	--ca-file cert/ca.crt 
	data binary get
#get specific cred
./client --workspace /path/to/workspace \
	--ca-file cert/ca.crt 
	data binary get --id 1
#delete specific cred 
./client --workspace /path/to/workspace \
	--ca-file cert/ca.crt 
	data binary delete --id 1
```

#### Sync
remove all data from local storage and write data from server
```bash
./client --workspace /path/to/workspace \
	--ca-file cert/ca.crt 
	data sync
```


Main schema.

When client execute Login operation, he send own public key to server. And server create relation public key with user id and session id and return JWT token to client.

Saving private data process
1) Server has master key or Key Encryption Key (KEK)
2) Generate new Data Encryption Key(DEK)
3) Encode private data with help DEK
4) Encode DEK with help KEK
5) Save it on storage
6) Encore DEK with help client's public key
7) Return to client
8) Client save data on local stora. On this moment client have encrypted data and encrypted DEK. His private key is right key for decryption DEK, and DEK can decrypte data.

If client would lose his public and private key, he just need to execute Login operation and send new key to server and after execute "sync" operation