# pm

[Difrex's perl password manager](https://github.com/difrex/pm) rewritten in go

# install

From AUR
```
pacaur -S pm
```

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

# examples

## adding a new password

### store your own password

```
./pm -wn porn -l coolpornsite.com -u john -p coolpassword -c 'my favorite site!'
```

### let pm generate a password by omitting the -p flag

```
./pm -wn porn -l coolpornsite.com -u john -c 'my favorite site!'
```

### choose the length of a generated password with the -L flag

```
./pm -wn porn -l coolpornsite.com -u john -L 32 -c 'my favorite site!'
```

### adding a password in the interactive way
```
./pm -wI
```

## showing passwords

### show passwords in dmenu

```
./pm -m
```

### show passwords in rofi
```
./pm -R
```

### print all passwords

```
$ ./pm -sn all
id: 1
name: porn
resource: coolpornsite.com
username: john
comment: my favorite site!
group:
```

### print all passwords related to the group 'work'

```
./pm -sg work
```

### print passwords in a nice formatted table by using the -t flag

```
$ ./pm -stn all
id name resource         username comment           group
----------------------------------------------------------
1  porn coolpornsite.com john     my favorite site!
```

### find password by name and copy it in the clipboard

```
$ ./pm -sn porn
password was copied to the clipboard!
URL: coolpornsite.com
User: john
Group:
```

### copy password in the clipboard and follow the link in the browser

```
./pm -son porn
```

## removing a password

```
$ ./pm -ri 13
successfuly removed password with id 13
```
