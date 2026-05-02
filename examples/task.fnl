(local crypto (require :crypto))

(fn task [val]
  (-> val
      (crypto.md5)
      (print)))
