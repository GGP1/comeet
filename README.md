# Comeet

Comeet is a service that keeps track of events from multiple calendar services and accounts, and sends you notifications before they start.

This project was implemented in an attempt to integrate virtual events in a natural way into my daily routine and to stop thinking about not missing them.

Another motivating factor was the fact that most calendar services' notification options and their integration with other services is pretty limited or non-existent.

## Installation

<details>
	<summary>Pre-compiled binaries</summary>

Linux, macOS, BSD and Windows pre-compiled binaries can be found [here](https://github.com/GGP1/comeet/releases).
</details>

<details>
	<summary>Homebrew (Tap)</summary>

```
brew install GGP1/tap/comeet
```
</details>

<details>
	<summary>Scoop (Windows)</summary>

```
scoop bucket add GGP1 https://github.com/GGP1/scoop-bucket.git
scoop install GGP1/comeet
```
or 
```
scoop install https://raw.githubusercontent.com/GGP1/scoop-bucket/master/bucket/comeet.json
```
</details>

<details>
	<summary>Docker</summary>

```
docker run -d gastonpalomeque/comeet -e COMEET_CONFIG=/config.yml -v /host/config/path:/config.yml
```
</details>

<details>
	<summary>Compile from source</summary>

> Requires Go 1.17 to be installed

```
git clone https://github.com/GGP1/comeet
cd comeet
go build main.go
```
</details>

## Configuration

The path to the configuration file must be set in the `COMEET_CONFIG` environment variable. A sample configuration file can be found at [config_sample.yaml](config_sample.yml).

Comeet sets itself to run on startup the first time it's executed.

For setting up different calendar and messaging services please follow their steps under [/docs/config/](/docs/config/).

## Why not webhooks?

Webhooks provide an efficient way of consuming real-time third-party API event updates but I decided **not** to use them for the following reasons:

- For this particular use case, polling works just well.
- Requires an instance in a cloud provider or a reverse-proxy, which are either paid or time-limited in free versions.
- Increased time-to-use and further configuration/maintenance.
- In the case of a reverse-proxy on localhost, security concerns due to the exposure.

However, this doesn't mean they aren't going to be an option in the future as many people may find them necessary.

## Caveats

- Updates during an event or within 15 minutes of its start won't be noticed by comeet, it's assumed that you will be aware of such changes due to the event's proximity.