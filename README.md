autodock
========

Perform actions based on a webhook from the Docker hub.

Idea: Make a simple Docker container that accepts a webhook from the Docker hub and triggers an action.

On successful build:

* `docker pull configured/image`
* `docker pull configured/image && restart_stuff.sh`

