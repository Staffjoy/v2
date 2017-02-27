# Crypto 

The crypto package standardizes "secure" things in Staffjoy.

## Hashing

To hash a password, use the built in hash lib.

To make a secret for storing in a db:
```
storeableSecret, err := crypto.HashPassword([]byte("VERY SECRET")
```

You can then verify a password attempt with:

````
err := crypto.CheckPassword(storeableSecret, []byte("PasswordPlaintext"))
if err != nil {
    // wrong password!
} else  {
    // Access granted
}
```

## Generating UUIDs

In general, don't use SQL auto-incremented integers as IDs. These are vulnerable
to enumeration attacks - so, a bad actor who gets access to one bad page can
keep increasing the id by 1 and finding secret information. In addition,
they reveal secret information - like how many users we have!

Instead, use a UUID. UUIDs are Universally Unique IDentifiers, and are standardized
by RFC4122.

https://en.wikipedia.org/wiki/Universally_unique_identifier

To generate a new UUID that is basically guaranteed to be unique across the internet:

```
uuid, err := crypto.NewUUID()
if err != nil {
    panic()
}

fmt.Printf("Your new UUID is %s", uuid)
```
