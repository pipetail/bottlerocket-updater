# Bottlerocket updater

[![.github/workflows/main-checker.yml](https://github.com/pipetail/bottlerocket-updater/actions/workflows/main-checker.yml/badge.svg)](https://github.com/pipetail/bottlerocket-updater/actions/workflows/main-checker.yml)

Bottlerocker updater is a set of two small executables
that help with the management of the OS updates in
your Kubernetes cluster.

## Updater

Updater should be running as daemon set on each of your
nodes. Updater communicates with Bottlerocket API
via the Unix Domain Socket and it does the following
sequence of activities:

1. `/actions/refresh-updates`
2. `/actions/prepare-update`
3. `/actions/activate-update`

## Checker

Checker is part that needs to be integrated to
[Kured](https://github.com/weaveworks/kured)
container image.

Then, it can be used with the Kured flag.

```
--reboot-sentinel-command string      command for which a successful run signals need to reboot (default ""). If non-empty, sentinel file will be ignored.
```

This utility basically just checks the update status
and it exits successfully when `update_state` is `Ready`.

Kured's DaemonSet manifest needs to altered this way:

```yaml
      volumes:
        - name: bottlerocket
          hostPath:
            path: /var/run/api.sock
      containers:
        - name: kured
          image: docker.io/weaveworks/kured:1.6.1
                 # If you find yourself here wondering why there is no
                 # :latest tag on Docker Hub,see the FAQ in the README
          imagePullPolicy: IfNotPresent
          securityContext:
            privileged: true # Give permission to nsenter /proc/1/ns/mnt
          volumeMounts:
            - name: bottlerocket
              mountPath: /var/run/bottlerocket.sock 
```

