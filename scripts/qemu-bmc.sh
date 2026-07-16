#!/bin/sh

set -eux

URL=https://github.com/legoater/qemu-aspeed-boot/raw/refs/heads/master/images/catalina-bmc/openbmc-20251217025154/obmc-phosphor-image-catalina-20251217025154.static.mtd.xz

IMAGE=obmc-phosphor-image-catalina-20251217025154.static.mtd
MACHINE=catalina-bmc
HOSTNAME=qemu-catalina

if [[ ! -e "${IMAGE}" ]]; then
  echo "downloading image"
  wget ${URL}

  echo "extracting image"
  unxz ${IMAGE}.xz
fi

qemu-system-arm -machine help

echo "running qemu"

qemu-system-arm \
 -m 1024 \
 -M "${MACHINE}" \
 -display none \
 -drive file="${IMAGE}",format=raw,if=mtd \
 -serial pty \
 -monitor unix:/tmp/qemu-"${MACHINE}"-hmp.sock,server,nowait \
 -qmp unix:/tmp/qemu-"${MACHINE}"-qmp.sock,server,nowait \
 -d guest_errors,unimp \
 -D /tmp/qemu.log \
 -nic user,hostfwd=tcp:0.0.0.0:2202-:22,hostfwd=tcp:0.0.0.0:8880-:80,hostfwd=tcp:0.0.0.0:4403-:443,hostname="${HOSTNAME}"
