# Identity Server

The Identity Server is a crucial component of our system, serving as the central hub for authentication via OAuth2
providers. It is designed to manage and authenticate users, ensuring secure access to our services.

## User Management

Users are the backbone of any system, and our Identity Server is no exception. It stores user data securely, ensuring
that each user's information is kept private and safe.

## OAuth2 Authentication

Our Identity Server uses OAuth2, an industry-standard protocol for authorization, to authenticate users. OAuth2 provides
a secure and reliable method for users to grant our services access to their information without sharing their password.

## Daemon/Robot Users

In addition to regular users, our Identity Server also manages daemon or robot users. These are non-human users that
perform automated tasks within our system. They are authenticated via API keys, providing them with the necessary
permissions to carry out their tasks without compromising the security of the system.

## Run Actions locally

### Run Tests

```bash
act -r --rm -W .github/workflows/ci.yml
```

### Build Docker image

```bash
act -r --rm -W .github/workflows/build.yml -j 'build' --secret-file ~/.energimind/dhub-deployer.env
```
