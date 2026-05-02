(local crypto (require :crypto))

(fn task [val]
  (-> (crypto.md5 val)
      (print)))
