# DbunderFS
[![Build Status](https://app.travis-ci.com/kos-v/dbunderfs.svg?branch=dev)](https://app.travis-ci.com/github/kos-v/dbunderfs)
[![Coverage Status](https://codecov.io/gh/kos-v/dbunderfs/branch/dev/graph/badge.svg)](https://codecov.io/gh/kos-v/dbunderfs)

#### Status
In development of the first version

#### Building
`git clone https://github.com/kos-v/dbunderfs.git`
`cd dbunderfs`
`make build`

#### Installation
`cp dbfs /usr/local/bin`
`dbfs migrate up "mysql://user:pass@127.0.0.1/db"`

#### Mount/Unmount
Mount: `dbfs mount "mysql://user:pass@127.0.0.1/db" /home/user/mount_point`

Unmount: `dbfs unmount /home/user/mount_point`
