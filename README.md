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

## Checker (Kured + custom addons)

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
  - name: bottlerocket-usr-bin
    persistentVolumeClaim:
      claimName: bottlerocket-bin
  initContainers:
    - name: install
      image: ghcr.io/pipetail/bottlerocket-updater/checker:linux-amd64-290643579b61c65824ff545a464af4da97f4904f
      volumeMounts:
        - mountPath: /bottlerocket/
          name: bottlerocket-usr-bin
      command:
        - /bin/bash
        - -c
        - |
          set -Eeuo pipefail
          mkdir -p /bottlerocket/bottlerocket
          cd /bottlerocket/bottlerocket
          cp /usr/local/bin/bottlerocket-checker ./
          cp /usr/local/bin/bottlerocket-reboot ./
          ls -lah ./
  containers:
    - name: kured
      image: ghcr.io/pipetail/bottlerocket-updater/checker:linux-amd64-290643579b61c65824ff545a464af4da97f4904f
      imagePullPolicy: IfNotPresent
      securityContext:
        privileged: true # Give permission to nsenter /proc/1/ns/mnt
      volumeMounts:
        - name: bottlerocket
          mountPath: /var/run/bottlerocket.sock
```

and

```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  labels:
    type: local
  name: bottlerocket-bin
  namespace: kube-system
spec:
  accessModes:
    - ReadWriteOnce
  capacity:
    storage: 2Gi
  hostPath:
    path: /opt/
    type: Directory
  storageClassName: manual
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: bottlerocket-bin
  namespace: kube-system
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
  storageClassName: manual
```

## Storage class

Kured operates with nsenter directly in the host operating
system. Hence InitContainer needs to copy check and reboot
binaries to the host filesystem.

This procedure is done via `manual` Storage Class that
allows access to host's filesystem.

```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: manual
provisioner: kubernetes.io/no-provisioner
volumeBindingMode: WaitForFirstConsumer
```