# Homie controller - WORK IN PROGRESS

Quick tool to collect informations from Homie-compliant sensors, and publish them in JSON format.

## Build

```
docker build -f  Dockerfile.build -q .
```

You can get the built binary by running

```
docker run --rm <id> /usr/local/bin/get-artifact > homie-controller
```

And then build a runtime image

```
docker build jbonachera/homie-controller .
```
