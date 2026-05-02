(local crypto (require :crypto))

(fn task [val]
  (print (crypto.md5 val)))
