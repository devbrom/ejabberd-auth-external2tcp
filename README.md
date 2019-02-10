# ejabberd-auth-external2tcp
A small external authentication script that translates data from /dev/stdin to stream and from stream to /dev/stdout

# ejabberd config example
```yaml
auth_method: [external]
extauth_program: "/path/to/compiled/script [addr]:[port]"
extauth_instances: 3
auth_use_cache: false
```
