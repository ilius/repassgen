# Summary
This program generates **fully customizable random passwords**, with almost **any desired pattern or length**, using the same syntax as [Regular Expressions](https://www.regular-expressions.info/quickstart.html) for **character classes** (POSIX-style) and **Repetition**.

See [examples](#examples) for clarity.

It's written in Go and does not have any external dependency (besides Go standard library)

# Build and install
Run:
```sh
go get
go install
```
This will compile and then copy `repassgen` executable file to `$GOPATH/bin/` directory.


# Feature Check List
- [x] Simple Repetition: `{N}`
- [x] Range Repetition: `{M-N}`
- [x] `[:alnum:]` 	Alphanumeric characters
- [x] `[:alpha:]` 	Alphabetic characters
- [x] `[:word:]` 	Word characters (letters, numbers and underscores)
- [x] `[:upper:]` 	Uppercase letters
- [x] `[:lower:]` 	Lowercase letters
- [x] `[:graph:]` 	Visible characters
- [x] `[:print:]` 	Visible characters and spaces (anything except control characters)
- [x] `[:digit:]` 	Digits
- [x] `[:xdigit:]` 	Hexadecimal digits
- [x] `[:punct:]` 	Punctuation (and symbols).
- [x] `[:space:]` 	All whitespace characters, including line breaks 
- [x] `[:blank:]` 	Space and tab
- [x] `[:cntrl:]` 	Control characters
- [x] `[:ascii:]` 	ASCII characters
- [x] `$base64(...)` Base64 encode function
- [x] `$base64url(...)` URL-safe Base64 encode function
- [x] `$base32(...)` Crockford's Base32 encode function (lowercase)
- [x] `$BASE32(...)` Crockford's Base32 encode function (uppercase)
- [ ] `$hex(...)` Hex encode function (lowercase)
- [ ] `$HEX(...)` Hex encode function (uppercase)



# Examples
- Alphanumeric password with length 12
    ```sh
    $ repassgen '[:alnum:]{12}'
    q8nrqhPQFNqj
    ```

- Alphabetic password with length 12
    ```sh
    $ repassgen '[:alpha:]{12}'
    wiADcFkhpjsk
    ```

- Lowercase alphabetic password with length 16
    ```sh
    $ repassgen '[:lower:]{16}'
    rnposblbuduotibe
    ```

- Uppercase alphabetic password with length 16
    ```sh
    $ repassgen '[:upper:]{16}'
    TQZZJHKQRKETOFNZ
    ```

- Numeric password with length 8
    ```sh
    $ repassgen '[:digit:]{8}'
    47036294
    ```

- A custom combination: 7 uppercase, space, 7 lowercase, space, 2 digits
    ```sh
    $ repassgen '[:upper:]{7} [:lower:]{7} [:digit:]{2}'
    UOHMGVM toubgvs 73
    ```

- Password with length 12, using only Base64 characters
    ```sh
    $ repassgen '[:b64:]{12}'
    6+BA71WCy90I
    ```

- Password with length 12, using only URL-safe Base64 characters
    ```sh
    $ repassgen '[:b64url:]{12}'
    j15w_qTncikR
    ```

- Password with length 16, using only Crockford's Base32 characters (lowercase)
    ```sh
    $ repassgen '[:b32:]{16}'
    zmt87n9hpcd2w28h
    ```

- Password with length 16, using only Crockford's Base32 characters (uppercase)
    ```sh
    $ repassgen '[:b32:]{16}'
    HJ48VSR4Y0DHQ9EV
    ```

- Password with length 16, using Crockford's Base32 characters and punctuations
    ```sh
    $ repassgen '[:b32::punct:]{16}'
    20s:z.5mbwws474y
    ```

- Specify character range manually: lowercase alphebetic with length 16
    ```sh
    $ repassgen '[a-z]{16}'
    qefqiocrabpiaags
    ```

- Specify character range manually: alphanumeric with length 12
    ```sh
    $ repassgen '[a-zA-Z0-9]{12}'
    XcwDAagzMlwv
    ```

- Include non-ramdom characters in the password (here: Test / , .)
    ```sh
    $ repassgen 'Test-[:alpha:]{4}/[:alpha:]{4},[:alpha:]{4}.[:alpha:]{4}'
    Test-Jcis/uLwq,SazR.CEFJ
    ```

- A 16-digit number similar to a credit card number
    ```sh
    repassgen '[:digit:]{4}-[:digit:]{4}-[:digit:]{4}-[:digit:]{4}'
    3996-9634-1459-0656
    ```

- Alphabetic password with a length 12 and 16 characters
    ```sh
    $ repassgen '[:alpha:]{12-16}'
    uamePKmuUUUcI
    ```

- Gerenate random bytes, then run Base64 encode function
    ```sh
    $ repassgen '$base64([:ascii:]{12})'
    YxRAUFhbFxRwSBxM
    ```

- Gerenate random bytes, then run Crockford's Base32 encode function
    ```sh
    $ repassgen '$base32([:ascii:]{12})'
    e274czjwcstd8v2ynv4
    ```
