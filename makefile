build:
	podman build -t harbor.pulze.cloud/home/abb-exporter:dev --arch=arm64,amd64 .
	podman push harbor.pulze.cloud/home/abb-exporter:dev