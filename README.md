# Sentry Sendmail

[![Build Status](https://travis-ci.org/retailnext/sentry-sendmail.svg?branch=master)](https://travis-ci.org/retailnext/sentry-sendmail)

This is sentry-sendmail, a sendmail replacement MTA. That means that all incoming emails on a system will get forwarded to sentry as a sentry issue.

This can be useful in situations where SMTP can not be used or local email server can not be installed, and we want to capture system's email as a sentry event.

## Installation

```bash
sudo apt-get install sentry-sendmail
```

## Local Build
Retrieve the latest copy of the source code by cloning the repository.

```bash
git clone https://github.com/retailnext/sentry-sendmail.git $HOME/sentry-sendmail
```

### Build packages using GoReleaser
Build sentry-sendmail from source in a single step using make.

```bash
goreleaser release --clean --snapshot
```

This will create installation tarballs and packages in `dist/`.

## Configuration

Enter the sentry DSN in the config file `/etc/sentry-sendmail.conf`