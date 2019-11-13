## Sentry Sendmail

This is sentry-sendmail, a sendmail replacement MTA. That means that all incoming emails on a system will get forwarded to sentry as a sentry issue.

This can be useful in situations where SMTP can not be used or local email server can not be installed, and we want to capture system's email as a sentry event.

## Installation

```bash
sudo apt-get install sentry-sendmail
```

## Manual Build
Retrieve the latest copy of the source code by cloning the repository.

```bash
git clone https://github.com/retailnext/sentry-sendmail.git $HOME/sentry-sendmail
```

## Make
Build sentry-sendmail from source in a single step using make.

```bash
cd $HOME/sentry-sendmail
make
```

Make will build the linux binary and the debian package in the `/dist` directory. And the binary for your platform on the `/bin` directory.

## Configuration

Enter the sentry DSN in the config file `/etc/sentry-sendmail.conf`