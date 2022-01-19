#!/bin/bash


setmqspl -m QMAPP2 -p APP2.APP1.FOO.RQ -s SHA1 -a "CN=bob,O=IBM,C=GB" -e AES256 -r "CN=alice,O=IBM,C=GB"
setmqspl -m QMAPP2 -p APP1.APP2.FOO.Q -s SHA1 -a "CN=alice,O=IBM,C=GB" -e AES256 -r "CN=bob,O=IBM,C=GB"


setmqaut -m QMAPP2 -t queue -n SYSTEM.PROTECTION.ERROR.QUEUE -p app +put
setmqaut -m QMAPP2 -t queue -n SYSTEM.PROTECTION.POLICY.QUEUE -p alice +browse

dspmqspl -m QMAPP2

bash

