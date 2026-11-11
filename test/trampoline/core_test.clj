(ns trampoline.core-test
  (:require [babashka.fs :as fs]
            [clojure.test :refer [deftest is testing]]
            [trampoline.core :as tr]))

(def home (tr/home))

(deftest expand-test
  (testing "a leading `~`"
    (is (= home (tr/expand "~")))
    (is (= (str home "/a") (tr/expand "~/a")))
    (is (= "~a" (tr/expand "~a")))
    (is (= "a/~" (tr/expand "a/~"))))
  (testing "env variables"
    (is (= (str home "/a") (tr/expand "$HOME/a")))
    (is (= (str home "/a") (tr/expand "${HOME}/a")))
    (is (= "/a" (tr/expand "$TRAMPOLINE_TEST_UNSET/a")))))

(deftest shorten-test
  (is (= "~/a" (tr/shorten (str home "/a"))))
  (is (= "/elsewhere" (tr/shorten "/elsewhere"))))

(deftest quote-zsh-test
  (is (= "/a/b-c" (tr/quote-zsh "/a/b-c")))
  (is (= "'/a b'" (tr/quote-zsh "/a b")))
  (is (= "'/it'\\''s'" (tr/quote-zsh "/it's"))))

(deftest quote-pwsh-test
  (is (= "C:\\a\\b" (tr/quote-pwsh "C:\\a\\b")))
  (is (= "'C:\\a b'" (tr/quote-pwsh "C:\\a b")))
  (is (= "'it''s'" (tr/quote-pwsh "it's"))))

(defn- with-tree
  "Calls `f` with a temp dir holding the dirs `a`, `.b`, `C` and `with space`, and the file
  `f.txt`."
  [f]
  (let [dir (fs/canonicalize (fs/create-temp-dir))]
    (try
      (doseq [d ["a" ".b" "C" "with space"]]
        (fs/create-dir (fs/path dir d)))
      (spit (str (fs/path dir "f.txt")) "")
      (f (str dir))
      (finally
        (fs/delete-tree dir)))))

(deftest sub-dirs-test
  (with-tree
    (fn [dir]
      (testing "only dirs, hidden ones included, sorted ignoring case and punctuation"
        (is (= (map #(str (fs/path dir %)) ["a" ".b" "C" "with space"])
               (tr/sub-dirs dir)))))))

(deftest items-test
  (with-tree
    (fn [dir]
      (let [file (str (fs/path dir "f.txt"))]
        (testing "a dir lists its sub-dirs"
          (is (= (map #(str (fs/path dir %)) ["a" ".b" "C" "with space"])
                 (map second (tr/items [dir])))))
        (testing "a dir with `:0` lists itself, displayed as written"
          (is (= [[dir dir]] (tr/items [(str dir ":0")]))))
        (testing "a file lists itself"
          (is (= [[file file]] (tr/items [file]))))
        (testing "a missing path lists nothing"
          (is (empty? (tr/items [(str (fs/path dir "missing"))]))))
        (testing "lines are listed in order"
          (is (= [file dir] (map second (tr/items [file (str dir ":0")])))))))))
