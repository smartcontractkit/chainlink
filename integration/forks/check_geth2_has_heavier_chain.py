#/usr/bin/python3

import os

height_cmd = '''docker-compose --compatibility logs geth%s | \
  egrep -o 'Commit new mining work +number=[0-9]+ ' | \
  sed 's/[^0-9]//g' |
  tail -1'''

geth_one_height = int(os.popen(height_cmd % 1).read())
geth_two_height = int(os.popen(height_cmd % 2).read())
assert geth_two_height > geth_one_height
