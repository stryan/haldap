# Home Assistant LDAP authenticator
A little program to interface with Home Assistant's CLI authentication provider system since the documented method no longer works in docker.

This program is designed to be either included in a customer HA docker build or bind mounted into the container.

## Building

For x86 systems:
	`go build`
For use in the Alpine-based HA docker container:
	`CGO_ENABLED=0 GOOS=linux go build -a`

## Home Assistant Config:
Add the following section to your `configuration.yml`:

```
homeassistant:
  auth_providers:
    - type: command_line
      command: /root/haldap
      args: ["--url=YOUR_LDAP_SERVER_HERE","--basedn=YOUR_LDAP_USER_BASE_DN_HERE","--userattr=YOUR_LDAP_USER_ATTRIBUTE_HERE","--admin=YOUR_LDAP_ADMIN_GROUP","--user=YOUR_LDAP_USER_GROUP"]
      meta: true
```

## LDAP config
This assumes you have `memberOf` enabled for groups and that you use the `displayName` attribute
