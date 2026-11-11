(ns build
  (:require [clojure.tools.build.api :as b]))

(def basis (delay (b/create-basis {:project "deps.edn"})))

(defn uber [_]
  (b/delete {:path "target"})
  (b/compile-clj {:basis @basis
                  :ns-compile '[trampoline.core]
                  :class-dir "target/classes"})
  (b/uber {:basis @basis
           :class-dir "target/classes"
           :uber-file "target/trampoline.jar"
           :main 'trampoline.core}))
