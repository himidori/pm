# pm

[Difrex's perl password manager](https://github.com/difrex/pm) rewritten in go

# usage

## gpg key 

generate a gpg key if you don't have one

```
gpg --gen-key
```

set your gpg key as the default one in ~/.gnupg/gpg.conf

or use a custom key instead 

```
# Create file with key email
cat > ~/.PM/.key << EOF
key_email@example.com
EOF
```

## first run 

```
$ ./pm
creating configuration directory...
creating database scheme...
encrypting database...
```

## check help

```
$ ./pm -h
Simple password manager written in Go

-s                      show password
-n [Name of resource]   name of resource
-g [Group name]         group name
-o                      open link
-t                      show passwords as table
-w                      store new password
-I                      interactive mode for adding new password
-l [Link]               link to resource
-u                      username
-c                      comment
-p [Password]           password
                        (if password is omitted PM will
                        generate a secure password)
-L [Length]             length of generated password
-r                      remove password
-i                      password ID
-m                      show dmenu
-R                      show rofi
-h                      show help
```
