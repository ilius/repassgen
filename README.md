# Summary

This program generates **fully customizable random passwords**, with almost **any desired pattern or length**, using the same syntax as [Regular Expressions](https://www.regular-expressions.info/quickstart.html) for **character classes** (POSIX-style), **repetition** (`{N}` or `{M,N}`) and **capturing groups** (`(...)`).

See [examples](#examples) for clarity.

[![Go Report Card](https://goreportcard.com/badge/github.com/ilius/repassgen)](https://goreportcard.com/report/github.com/ilius/repassgen)

# Build and install

Run:

```sh
go get github.com/ilius/repassgen
```

This will compile and then copy `repassgen` executable file to `$GOPATH/bin/` directory.

If you just want to compile without installing it

```sh
go get -d github.com/ilius/repassgen
cd $GOPATH/src/github.com/ilius/repassgen
go build
```

Then `repassgen` binary file will be created in current directory.

# Features of regexp

- \[x\] Simple repetition: `{N}`
- \[x\] Range repetition: `{M,N}`
- \[x\] Manual character range, like `[a-z1-579]`
- \[x\] Repeatable groups with `(...){N}`, like  `([a-z]{2}[1-9]){3}`
- \[x\] `[:alnum:]` Alphanumeric characters
- \[x\] `[:alpha:]` Alphabetic characters
- \[x\] `[:word:]`, `[:w:]` or `\w`: Word characters (letters, numbers and underscores, same as `[a-zA-Z0-9_]`)
- \[x\] `[:upper:]` Uppercase letters
- \[x\] `[:lower:]` Lowercase letters
- \[x\] `[:graph:]` Visible characters
- \[x\] `[:print:]` Visible characters and spaces (anything except control characters)
- \[x\] `[:digit:]` or `\d` Digits
- \[x\] `[:xdigit:]` Hexadecimal digits
- \[x\] `[:punct:]` Punctuation (and symbols).
- \[x\] `[:space:]` All whitespace characters, including line breaks
- \[x\] `[:blank:]` Space and tab
- \[x\] `[:cntrl:]` Control characters
- \[x\] `[:ascii:]` ASCII characters
- \[x\] [Unicode code points](https://www.regular-expressions.info/unicode.html), like `[\u00e0-\u00ef]{5}`
- \[x\] Group references `\1`, `\2`, etc

# Aditional Features (not part of regexp)

- \[x\] Combined multiple named/manual character classes, for example:
  - `[:digit:a-m]`
  - `[:digit::alpha:]` = `[:alnum:]`
- \[x\] `[:b32:]` Crockford's Base32 alphabet (lowercase)
- \[x\] `[:B32:]` Crockford's Base32 alphabet (uppercase)
- \[x\] `[:B32STD:]` Standard Base32 alphabet (uppercase)
- \[x\] `[:b64:]` Standard Base64 alphabet
- \[x\] `[:b64url:]` URL-safe Base64 alphabet
- \[x\] `$base64(...)` Base64 encode function (input is hex-encoded)
- \[x\] `$base64url(...)` URL-safe Base64 encode function (input is hex-encoded)
- \[x\] `$base32(...)` Crockford's Base32 encode function (lowercase) (input is hex-encoded)
- \[x\] `$BASE32(...)` Crockford's Base32 encode function (uppercase) (input is hex-encoded)
- \[x\] `$base32std(...)` Standard Base32 encode function (uppercase, with no padding) (input is hex-encoded)
- \[x\] `$hex(...)` Hex encode function (lowercase)
- \[x\] `$HEX(...)` Hex encode function (uppercase)
- \[x\] Show [entropy](https://en.wikipedia.org/wiki/Password_strength#Entropy_as_a_measure_of_password_strength) of pattern
  - Use `repassgen -entropy 'PATTERN'` command
  - Indicates strength of generated passwords, the higher the better
  - We recommand at least 47 bits (equal to 8 alphanumeric: `[:alnum:]{8}`)
  - Entropy of pattern is more important than entropy of password, if you re-use patterns
- \[x\] `$hex2dec(...)` Convert hexadecimal number to decimal number
- \[x\] `$escape(...)` Escape unicode characters, non-printable characters and double quote
- \[x\] `$?(...)` Randomly include or omit the string/pattern (%50 chance, adds 1 bit to entropy)
- \[x\] `$bip39word(N)` Generate N words from BIP-39 English mnemonic words
- \[x\] `$bip39encode(...)` Encode hex-encoded bytes into some BIP-39 English mnemonic words
- \[x\] `$date(2000,2020,-)` Generate a random date in the given year range
- \[x\] `$space(...)` Adds spaces between each two characters of string (generated from given pattern)
- \[x\] `$expand(|...)` Adds `|` (for example) between each two characters (similar to `$space`)
- \[x\] `$rjust(PATTERN,N,X)` Justify to right, `N` is width (N>=1), `X` is the character to fill
- \[x\] `$ljust(PATTERN,N,X)` Justify to left, similar to `$rjust`
- \[x\] `$center(PATTERN,N,X)` Justify to center, similar to `$rjust`
- \[x\] `$pyhex(...)` Convert hex-encoded bytes to Python `bytes` with hex values (like `b'\x74\x65\x73\x74'`)
- \[x\] `$romaji(...)` Converts Japanese hiragana/katakana string to Latin

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
  $ repassgen '[:B32:]{16}'
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

- Include non-random characters in the password (here: Test / , .)

  ```sh
  $ repassgen 'Test-[:alpha:]{4}/[:alpha:]{4},[:alpha:]{4}.[:alpha:]{4}'
  Test-Jcis/uLwq,SazR.CEFJ
  ```

- A 16-digit number similar to a credit card number

  ```sh
  repassgen '[:digit:]{4}-[:digit:]{4}-[:digit:]{4}-[:digit:]{4}'
  3996-9634-1459-0656
  ```

- Alphabetic password with a length between 12 and 16 characters

  ```sh
  $ repassgen '[:alpha:]{12,16}'
  uamePKmuUUUcI
  ```

- Gerenate random bytes, then run Base64 encode function

  ```sh
  $ repassgen '$base64($byte(){12})'
  bsOuN8KuRsOFw5jClkDDjMOrFA==
  ```

- Gerenate random bytes, then run Crockford's Base32 encode function

  ```sh
  $ repassgen '$base32($byte(){12})'
  c51e2kk1aafe3jngm3gxqrazpwqva
  ```

- Use `-` or `[` or `]` inside `[...]`

  ```sh
  $ repassgen '[.\- ]{50}'
  - .-.-- --.------- --.. -.---.-.. -- --..-..  .---
  ```

  ```sh
  $ repassgen '[a-z\[\]\-]{20}'
  nylhjcdq-qcaajvpcxuo
  ```

- Use whitespace characters like newline or tab (inside or outside `[...]`)

  ```sh
  $ repassgen '[a-z\t]{5}\t[a-z\t]{5}\n[a-z\t]{10}'
  caelk	zccqm
  zpbgjba	pm
  ```

- Generate 12 random mnemonic words from [BIP-39](https://en.bitcoin.it/wiki/BIP_0039) English words

  ```sh
  $ repassgen '$bip39word(12)'
  cinnamon purity funny pigeon arrive equal foil alter life accident bar roast
  ```

- Generate 16 random bytes, then encode it to [BIP-39](https://en.bitcoin.it/wiki/BIP_0039) English mnemonic words

  ```sh
  $ repassgen '$bip39encode($byte(){16})'
  useful come fall plunge breeze side skill another boil expose essence about
  ```
