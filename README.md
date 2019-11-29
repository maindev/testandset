# TestAndSet CLI

With the TestAndSet CLI you can handle custom web-based mutexes from the TestAndSet API. Mutexes can be integrated in your scripts, whereever you need them.

Commands
========

`lock` - You can create a mutex with your custom name. The output can be a json-object or a pure refresh token. You can set a timeout, in which `lock` automatically tries to lock again at intervals.

Possible parameters are: `--name, -n $MUTEX_NAME`, `--output, -o [json|token]` (optional), `--timeout, -t (time in seconds)` (optional) and `--owner, -O $MUTEX_OWNER_NAME`

```
$ testandset mutex lock --name MyMutex --output json --timeout 60 --owner main.dev
```

`get` - You can check if a mutex with a specific name currently exists.

Possible parameters are: `--name, -n $MUTEX_NAME`

```
$ testandset mutex get --name MyMutex
```

`refresh` - You can refresh an existing mutex for a small amount of time. To check that you are allowed to refresh you need the token provided by the `lock` command.

Possible parameters are: `--name, -n $MUTEX_NAME` and `--token, -t $TOKEN`

```
$ testandset mutex refresh --name MyMutex --token abcd1234-567e-890f-g123-hijk4567
```

`unlock` - You can unlock an existing mutex to free it for other users when your work is done. To check that you are allowed to unlock you need the token provided by the `lock` command.

Possible parameters are: `--name, -n $MUTEX_NAME` and `--token, -t $TOKEN`

```
$ testandset mutex unlock --name MyMutex --token abcd1234-567e-890f-g123-hijk4567
```

`auto-refresh` - You can automatically let the CLI refresh an existing mutex at intervals. When this command is stopped it will automatically unlock the mutex. To check that you are allowed to auto-refresh you need the token provided by the `lock` command.

Possible parameters are: `--name, -n $MUTEX_NAME` and `--token, -t $TOKEN`

```
$ testandset mutex auto-refresh --name MyMutex --token abcd1234-567e-890f-g123-hijk4567
```

Example
=======

Here is an example shell script how you can use TestAndSet CLI

```
#!/bin/bash -eu

MUTEX_NAME=maindev-mutext-example-1
TOKEN=$(./TestAndSet mutex lock --name $MUTEX_NAME --output token --timeout 120)
./TestAndSet mutex auto-refresh --name $MUTEX_NAME --token $TOKEN &
PID=$!
trap "kill -TERM $PID" EXIT

echo "Doing something important!"
sleep 20
echo "Done."
```

Licensing
=========
TestAndSet CLI is licensed under the New BSD License. See
[LICENSE](https://github.com/maindev/testandset/blob/master/LICENSE) for the full
license text.