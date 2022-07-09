# Init

Sample configuration files for:

```
systemd: broln.service
```

## systemd

Add the example `broln.service` file to `/etc/systemd/system/` and modify it according to your system and user configuration. Use the following commands to interact with the service:

```bash
# Enable broln to automatically start on system boot
systemctl enable broln

# Start broln
systemctl start broln

# Restart broln
systemctl restart broln

# Stop broln
systemctl stop broln
```

Systemd will attempt to restart broln automatically if it crashes or otherwise stops unexpectedly.
