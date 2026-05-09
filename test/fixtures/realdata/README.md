# Real-Data Fixtures

Most fixtures are copied from the Apache SpamAssassin public mail corpus:

https://spamassassin.apache.org/old/publiccorpus/readme.html

The corpus readme states that headers are reproduced in full, some address obfuscation occurred, and the messages were public forum posts, public newsletters, messages sent with knowledge they may be public, or messages sent by the corpus maintainer.

The calendar fixture is an RFC5545-style invite used because the public corpus slice does not include a representative `text/calendar` invite.

The truncated fixture is a partial derivative of a public corpus message and exists to test partial MIME recovery.

The huge fixture is a public corpus seed; tests expand it in memory rather than committing a multi-megabyte copy.

