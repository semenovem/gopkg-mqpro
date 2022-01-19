#!/bin/bash

setmqspl -m QMAPP1 -p APP1.APP2.FOO.RQ -s SHA1 -a "CN=alice,O=IBM,C=GB" -e AES256 -r "CN=bob,O=IBM,C=GB"
setmqspl -m QMAPP1 -p APP2.APP1.FOO.Q -s SHA1 -a "CN=bob,O=IBM,C=GB" -e AES256 -r "CN=alice,O=IBM,C=GB"

setmqaut -m QMAPP1 -t queue -n SYSTEM.PROTECTION.ERROR.QUEUE -p app +put
setmqaut -m QMAPP1 -t queue -n SYSTEM.PROTECTION.POLICY.QUEUE -p app +browse


dspmqspl -m QMAPP1

bash
#exit
#
##setmqaut -m QMAPP1 -t qmgr -p app +connect +inq
##setmqaut -m QMAPP1 -n TEST.Q -t queue -p app +put
##setmqaut -m QMAPP1 -n TEST.Q -t queue -p app +get
#
#
#setmqaut -m QMAPP1 -t queue -n SYSTEM.PROTECTION.POLICY.QUEUE -p app  +browse
#setmqaut -m QMAPP1 -t queue -n SYSTEM.PROTECTION.ERROR.QUEUE -p app  +put
#setmqspl -m QMAPP1 -p TEST.Q -s SHA1 -a "CN=alice,O=IBM,C=GB" -e AES256 -r "CN=bob,O=IBM,C=GB"
#
#dspmqspl -m QMAPP1
#
#
#exit
#setmqaut -m QMAPP1 -t qmgr -p app +connect +inq
#setmqaut -m QMAPP1 -n APP2.APP1.FOO.Q -t queue -p app +put
#setmqaut -m QMAPP1 -n APP2.APP1.FOO.Q -t queue -p app +get
#
#
#setmqaut -m QMAPP1 -t queue -n SYSTEM.PROTECTION.POLICY.QUEUE -p app +browse
#setmqaut -m QMAPP1 -t queue -n SYSTEM.PROTECTION.ERROR.QUEUE -p app +put
#
##setmqspl -m QMAPP1 -p TEST.Q -s SHA1 -a "CN=alice,O=IBM,C=GB" -e AES256 -r "CN=bob,O=IBM,C=GB"
#
#setmqspl -m QMAPP1 -p APP2.APP1.FOO.Q -s SHA512 -a "CN=alice,O=IBM,C=GB" -e NONE
