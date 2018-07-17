# pm

[Difrex's perl password manager](https://github.com/difrex/pm) rewritten in go

# usage

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

