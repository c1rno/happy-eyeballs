# Happy eyeballs

**[Source](https://blog.softwaremill.com/happy-eyeballs-algorithm-using-zio-120997ba5152)**
or if previous link will die here is a backup [algorithm](https://tools.ietf.org/html/rfc8305)

It's an algorithm which receives addresses list as a parameters
 and connecting to it one by one, if after a bit delay connecting
 still in process (or return error) it algorithm concurrently starts
 connecting to next address.
 If one of previous connections completed it stop all and return
 actual connect.
