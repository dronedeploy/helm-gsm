# helm-gsm

Based on [helm-secrets](https://github.com/jkroepke/helm-secrets), but I needed [Google Secret Manager](https://cloud.google.com/secret-manager) instead of [sops](https://github.com/mozilla/sops/), and I didn't want to install [sops](https://github.com/mozilla/sops/) if I didn't need it.

## Dependencies & Platforms

No dependencies, but only supports darwin_amd64 and linux_amd64 because that's all I needed. It would be trivial to compile for additional platforms if needed.

## Installation

```
helm plugin install https://github.com/dronedeploy/helm-gsm.git --version v0.2.0
```

## Usage

Takes a prepared yaml, grabs the plaintext secret from Google Secret Manager, and creates a decrypted file with b64 encoded plaintext secrets. By default it looks for `secrets.yaml` in the current directory and outputs `secrets.yaml.dec` in the current directory.
```
helm gsm
or
helm gsm -f path/to/secrets.yaml
```

Or alternatively use inline as a protocol handler - eg.
```
helm template -f gsm://path/to/secrets.yaml .
```

The secrets file needs to have a very specific format:
```
secrets:
  first_secret: gsm:my-project/super_secret/1
  second_secret: gsm:my-other-project/other_secret/3
  ...
```

The keys can be any valid Helm value key, and the encrypted string can be any value that GSM accepts, but the format for the reference value needs to be specific:

```
Secret reference format:
gsm:project_id/secret_name/version

From Google's API Documentation:
project_id: "The unique, user-assigned ID of the Project. It must be 6 to 30 lowercase letters, digits, or hyphens. It must start with a letter. Trailing hyphens are prohibited." 
secret_name: "Secret names can only contain English letters (A-Z), numbers (0-9), dashes (-), and underscores (_)"
version: Versions are a monotonically increasing integer starting at 1.

Regex:
^gsm:[a-z][a-z0-9-]{4,28}[a-z0-9]\/[a-zA-Z0-9_-]+\/[1-9]?[0-9]+$
```

## Other Notes
It's recommended that you add `*.dec` to your `.gitignore` file.
