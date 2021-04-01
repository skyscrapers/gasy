# GO-AWS-STS-YubiKey (GASY) - Deprecated

:exclamation::exclamation: **DEPRECATED :exclamation::exclamation: - this project is deprecated and will be deleted in the near future. Most of its features can be found in other open source tools like [aws-vault](https://github.com/99designs/aws-vault).**

A CLI tool to create AWS STS credentials and URLs using a YubiKey as your OTP.

## Usage

```bash
gasy -h
A CLI tool to generate STS keys and URLs using Yubikey OTP.

Please see the README for documentation: https://github.com/skyscrapers/gasy

Usage:
  gasy [flags]
  gasy [command]

Available Commands:
  accounts    List all available accounts
  help        Help about any command

Flags:
  -c, --client-list-location string   Path to the json client list
      --config string                 config file (default is $HOME/.gasy.toml)
  -h, --help                          help for gasy
  -p, --profile string                which AWS profile to use to perform the login (default "default")
  -r, --region string                 region to use with AWS (default "eu-west-1")
  -s, --serialnumber string           serial number of your AWS MFA device
  -S, --slotname string               Name of your YubiKey ath slot

Use "gasy [command] --help" for more information about a command.
```

Gasy expects 1 argument, the account number. You can view a list of available accounts by running `gasy accounts`

### Example

```bash
gasy 1
Using config file: /Users/dev/.gasy.toml
Please touch your YubiKey...
requesting credentials for Dev-account
Credentials written to profile!

export AWS_PROFILE=Dev-account

URL: https://signin.aws.amazon.com/federation?Action=login&...
```

```bash
gasy accounts
Using config file: /Users/dev/.gasy.toml
+----+-------------------------+--------------------------------+--------------+
| #  |          NAME           |          DESCRIPTION           |      ID      |
+----+-------------------------+--------------------------------+--------------+
|  0 | account1                | Main account                   | 123456789012 |
|  1 | account2                | test account                   | 123456789098 |
+----+-------------------------+--------------------------------+--------------+
```

## Configuration

gasy can be configured by passing the required flags (see `gasy -h`).
You can also write some of the parameters in a configuration file:
`~/.gasy.toml`:

```toml
[aws]
clientListLocation = "/Users/dev/aws-accounts.json"
mfaSerial = "arn:aws:iam::123456789012:mfa/user"

[yubikey]
slotName = "Amazon Web Services:user@account"
```

The clientlist is expected to be a json file using the following format:

```json
{
  "accounts": [
    {
      "id": "123456789012",
      "name": "account1",
      "sid": "91bb981a-12if-475a-a940-fc67abcddf10",
      "description": "Main account"
    },
    {
      "id": "123456789098",
      "name": "account2",
      "sid": null,
      "description": "test account"
    }
  ]
}

```

## Dependencies

Gasy depends on the `ykman` CLI tool provided by YubiCo.
It can be installed on mac with brew: `brew install ykman`

## TODO

- tests
- CI integration
- error handling
- ...
