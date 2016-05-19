# Oppositus -- mirror CoreOS releases

Package oppositus maintains a local mirror of CoreOS releases. _Oppositus_
is when your core is mirrored -- that is, your internal organs are reversed.

Only files with good signatures are stored, so the mirror can be safely used
via the local filesystem without MITM risk.

There is **no** incremental downloading; the update mechanism in
CoreOS is way too elaborate to imitate easily. I am also personally
most interested in the PXE image
[for running VMs](http://eagain.net/blog/2016/03/31/coreos-in-a-vm.html),
which is outside of that scope anyway.


## Usage

Config file sets which release channels to mirror (default: all), and
what files to include (default: all). First matching filter applies,
`-` excludes and `+` includes a file.


```console
$ cat config.json
{
    "channels": ["stable", "beta"],
    "filters": [
        "+ coreos_production_pxe[._]*",
        "+ coreos_developer_container[._]*",
        "- *"
    ]
}
$ mkdir dest
$ oppositus config.json dest
...
$ tree dest
dest
├── all
│   ├── 1010.3.0
│   │   ├── coreos_developer_container.bin.bz2
│   │   ├── coreos_developer_container.bin.bz2.DIGESTS
│   │   ├── coreos_developer_container.bin.bz2.DIGESTS.sig
│   │   ├── coreos_developer_container.bin.bz2.sig
...
│   │   ├── coreos_production_pxe_image.cpio.gz
│   │   ├── coreos_production_pxe_image.cpio.gz.sig
│   │   ├── coreos_production_pxe.README
│   │   ├── coreos_production_pxe.README.sig
│   │   ├── coreos_production_pxe.sh
│   │   ├── coreos_production_pxe.sh.sig
│   │   ├── coreos_production_pxe.vmlinuz
│   │   └── coreos_production_pxe.vmlinuz.sig
│   └── 899.17.0
│       ├── coreos_developer_container.bin.bz2
│       ├── coreos_developer_container.bin.bz2.DIGESTS
│       ├── coreos_developer_container.bin.bz2.DIGESTS.sig
│       ├── coreos_developer_container.bin.bz2.sig
...
│       ├── coreos_production_pxe_image.cpio.gz
│       ├── coreos_production_pxe_image.cpio.gz.sig
│       ├── coreos_production_pxe.README
│       ├── coreos_production_pxe.README.sig
│       ├── coreos_production_pxe.sh
│       ├── coreos_production_pxe.sh.sig
│       ├── coreos_production_pxe.vmlinuz
│       └── coreos_production_pxe.vmlinuz.sig
├── beta
│   └── current -> ../all/1010.3.0
└── stable
    └── current -> ../all/899.17.0

7 directories, 40 files
$ head -3 dest/stable/current/coreos_production_pxe.README
If you have qemu installed (or in the SDK), you can start the image with:
  cd path/to/image
  ./coreos_production_pxe.sh -curses
```

## TODO

- container to run it, systemd timer to schedule it
- garbage collection
- perhaps maintain symlinks in `<channel>/<version>` to note that said
  version was seen in that channel at some point in time
- use readOnlyRootFS in container manifest
