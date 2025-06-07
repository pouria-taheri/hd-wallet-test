### generate template config
```shell script
./hdwallet generate
```
or
```shell script
go run main.go generate
```

This command will generate template file containing random mnemonic & seed information(with default 256 entropy bit size).

> to change entropy bit size pass argument as entropy bit size between 128-256 and multiple of 32
> e.g. `./hdwallet generate 128`. this will generate template data containing random mnemonic & seed with 128 bit size entropy.

> to change template file path pass file name as second argument:
> e.g. `./hdwallet generate 128 template-config.yaml`
---

After generation of sample config, you should edit fields existing in file so 
application could load the file with correct settings.

```yaml
mnemonic: prosper oblige nice quantum effort giraffe donor blossom demise nurse vast forget
seed: 
password:
```
> only one of fields of mnemonic or seed will be required for application. 
> in the first step, application checks for `mnemonic`, if field exists in 
> config, `seed` will be ignored.

after configuring the config file, you must set the config path in main config
file(`config.yaml`) so application could recognise the secured config file.  

### lock existing config
After you set the config, application is able to run with unencrypted file.
To encrypt or change password of file, run:
```shell script
./hdwallet encrypt
```
or
```shell script
go run main.go encrypt
```
This command will, ask for new passwords to encrypt the file with.  
 
If the file has already been encrypted before, before settings new password, 
application asks for the current password of config to decrypt the encrypted file
in the first place.


The existing config file will be replaced by the new encrypted file.

### unlock config
If you want to set config as plain text again for further configurations, run:
```shell script
./hdwallet decrypt
```
or
```shell script
go run main.go decrypt
```
This command will ask for the current password of file if it's encrypted.

---

In `config.yaml` there is a field named `salt`, which is being used for generation of
encryption/decryption key combined with password coming from user input in CLI.

---

### GRPC

to handle rpc, client configuration must be enabled:

```yaml
salt: c309278c44aa43f9bd45c6e9bc203c86
cmd:
  server:
    grpc:
      address: 0.0.0.0:5050
      enabled: false
      tlsCertFile: cert/server-cert.pem
      tlsKeyFile: cert/cert/server-key.pem
  client:
    cli:
      enabled: true
    grpc:
      address: 127.0.0.1:5050
      enabled: false
      insecure: true
      CAFile: cert/ca-cert.pem
```

- `salt` is server key which will be combined with client password to encrypt/decrypt config
- `cmd.server.grpc.address` is address of gRPC server for remote unlock connections
- `cmd.server.grpc.address` if set true, gRPC server will be served on project start, otherwise
                            will be ignored, and CLI input will be shown for file decryption
- `cmd.server.grpc.tlsCertFile` optional TLS certificate file
- `cmd.server.grpc.tlsKeyFile` optional TLS key file
- `cmd.client.cli.enabled` used for turning off/on cli handler. turn cli off of if you're running
                           application in non-interactive terminal.
- `cmd.client.grpc.address` is address of remote server to connect
- `cmd.client.grpc.CAFile` optional is transport credentials file for requesting to the server

you can get list of available commands of application by passing param -h

```shell script
go run main.go -h
```
or
```shell script
./hdwallet -h
```