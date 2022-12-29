## How to launch with docker

1. Download config.yaml from this repository
2. Edit it
3. `docker run -v $(pwd)/config.yaml:/data/config.yaml -v /var/run/docker.sock:/var/run/docker.sock ghcr.io/ozraru/containerstarterbot`  