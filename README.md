autodock
========

Perform actions based on a webhook from [the Docker Hub](https://hub.docker.com/).

Idea: Make a simple Docker container that accepts a webhook from the Docker hub and triggers an action.

On successful build:

* `docker pull configured/image`
* `docker pull configured/image && restart_stuff.sh`

To do that, we'll need access to the host box.

## Quick run

```
docker run \
  --publish 8080:8080 \
  -e 'AUTODOCK_WEBAPP=training/webapp:echo yay' \
  autodock
```

Replace `training/webapp` with the name of the container from the Docker hub you want to trigger on.

## Development

Simply clone this repo, start hacking then run:

```
docker build -t autodock .
```

Make sure you run with some set of AUTODOCKs set up:

```
docker run -e 'AUTODOCK_YAY=example/app:echo hi' -p 8080:8080 autodock
```

You can then verify this with your favorite way to hit up a URL:

#### Python, with Requests

```python
>>> import json
>>> requests.post("http://localhost:8080/autodock/v1/",
...               data=json.dumps({"repository":{"repo_name":"example/app"}}),
...               headers={'Content-type': 'application/json'})
<Response [200]>
```

#### cURL
```
curl -X POST -H "Content-Type: application/json" \
     -d '{"repository":{"repo_name":"something/else"}}' \
     127.0.0.1:8080/autodock/v1/
```

### Example Run

```
$ docker build -t autodock .
Sending build context to Docker daemon 102.4 kB
Sending build context to Docker daemon
Step 0 : FROM google/golang-runtime
# Executing 2 build triggers
Step onbuild-0 : ADD . /gopath/src/app/
 ---> b51b592cfe60
Step onbuild-1 : RUN /bin/go-build
 ---> Running in 2d9377aecd4c
 ---> 88a51b3970fb
 ---> 88a51b3970fb
Removing intermediate container b81a681a3019
Removing intermediate container 2d9377aecd4c
Successfully built 88a51b3970fb
$ docker run -e 'AUTODOCK_YAY=example/app:echo hi' -p 8080:8080 autodock
2014/07/28 14:58:45 Docker repository actions:
2014/07/28 14:58:45 	example/app: [echo hi]
2014/07/28 14:58:57 Processing example/app
2014/07/28 14:58:57 Running [echo hi]
2014/07/28 14:58:57 hi
```



