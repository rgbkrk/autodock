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

