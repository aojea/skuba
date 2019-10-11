# Sandbox to create a development environment

Usually you need to test this or other functionality and you need an easy and fast way to do it.
There are several projects that are able to provide a kubernetes environment like kind of kubernetes, but 
they doesn't fulfill the requirements for skuba.

Skuba needs to provide the infrastructure, and ideally it should allow to use the SLES packages and kernel.

## MicroVMs to the rescue

[Firecracker](https://github.com/firecracker-microvm/firecracker) is a project that allows to spawn MicroVMs in subseconds.

There is a tool called [firectl](https://github.com/firecracker-microvm/firectl) that is a basic command-line tool that lets you run arbitrary Firecracker MicroVMs via the command line. 

This lets you run a fully functional Firecracker MicroVM, including console access, read/write access to filesystems, and network connectivi

### Install firecracker and firectl

Install firecracker latest release, at this times is version 0.18:

```sh
wget https://github.com/firecracker-microvm/firecracker/releases/download/v0.18.0/firecracker-v0.18.0 -o firecracker
chmod +x firecracker
mv firecracker /usr/local/bin/firecracker
```

Install firectl:



### Create the base images

We will need a root filesystem and a kernel image https://github.com/firecracker-microvm/firecracker/blob/master/docs/rootfs-and-kernel-setup.md                              
Obtain SUSE distro rootfs, per example:
                                                                                                                             
```
wget https://download.opensuse.org/repositories/Cloud:/Images:/Leap_15.1/images/openSUSE-Leap-15.1-OpenStack-rootfs.x86_64-0.0.4-Build7.44.tar.xz
```

Obtain the kernel image we want to test (firecracker doesn't have initrd so the main modules have to be compiled in the kernel image):


```sh
wget https://download.opensuse.org/repositories/Kernel:/HEAD/standard/x86_64/kernel-kvmsmall-5.3.rc8-2.1.gd6f0b71.x86_64.rpm    
# Obtain the kernel uncompressed image
rpm2cpio kernel-kvmsmall-5.3.rc8-2.1.gd6f0b71.x86_64.rpm | cpio -id                                                 
cd boot/
gunzip vmlinux-5.3.0-rc8-2.gd6f0b71-kvmsmall.gz
mv vmlinux-5.3.0-rc8-2.gd6f0b71-kvmsmall vmlinux
```

### Create the network

You have to create as many interfaces as VMs:

```sh
sudo ip tuntap add tap0 mode tap
sudo ip tuntap add tap1 mode tap
[...]
```

You have to create a linux bridge and provide dhcp addresses:




### Create the VMs

firectl --kernel=hello-vmlinux.bin --root-drive=alpine.ext4
--firecracker-binary=./firecracker --tap-device=tap0/aa:bb:11:22:33:44                                                      
--tap-device=tap1/aa:bb:11:22:33:11


## Create the kubernetes cluster with skuba

### Install skuba


### Install kubernetes

## References

[1] https://github.com/firecracker-microvm/firecracker/blob/master/docs/getting-started.md                                      