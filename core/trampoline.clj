#!/usr/bin/env bb

;; Shell-independent core of trampoline. The shell layers (trampoline.zsh, pwsh/trampoline.ps1)
;; only bind a key, run this script and execute the command it prints.
;;
;; Usage:
;;   bb trampoline.clj items [--type <d|f>]  Print the configured items, one per line, optionally
;;                                           only the directories (d) or the files (f).
;;   bb trampoline.clj pick --shell <shell>  Pick an item in fzf and print the command to run it,
;;                                           quoted for <shell>: zsh or pwsh. Prints nothing when
;;                                           the pick is cancelled.

(require '[babashka.fs :as fs]
         '[babashka.process :as p]
         '[clojure.string :as str])

#_[:type Shell [:enum "zsh" "pwsh"]]

;; `[display path]`. `display` keeps the line as written in the config (e.g. `${HISTFILE}`), while
;; `path` is the expanded filepath.
#_[:type Item [:tuple :string :string]]

;; `$HOME` first, as the shells use it for `~`. Java's home (`fs/home`) ignores it.
(def home #_[:- :string]
  (or (System/getenv "HOME") (str (fs/home))))

(def config-files #_[:- [:vector :string]]
  [(str (fs/path home ".trampoline")) (str (fs/path home ".trampoline.local"))])

(defn expand
  "Expands a leading `~` and the environment variables `$VAR` and `${VAR}`."
  #_[:=> [:cat :string] :string]
  [s]
  (-> s
      (str/replace #"^~(?=$|[/\\])" (str/re-quote-replacement home))
      (str/replace #"\$\{(\w+)\}|\$(\w+)"
                   (fn [[_ braced bare]]
                     (str/re-quote-replacement (or (System/getenv (or braced bare)) ""))))))

(defn shorten
  "Replaces the home directory prefix with `~`, for display."
  #_[:=> [:cat :string] :string]
  [s]
  (if (str/starts-with? s home)
    (str "~" (subs s (count home)))
    s))

(defn config-lines
  #_[:=> [:cat] [:sequential :string]]
  []
  (->> config-files
       (filter fs/regular-file?)
       (mapcat fs/read-all-lines)
       (remove str/blank?)))

(defn sub-dirs
  "Level-1 sub-dirs, hidden ones included. Sorted ignoring case and punctuation, like `sort` in an
  English locale: `a`, `.b`, `C`."
  #_[:=> [:cat :string] [:sequential :string]]
  [dir]
  (->> (fs/list-dir dir)
       (filter fs/directory?)
       (map (comp str fs/normalize))
       (sort-by (juxt #(str/lower-case (str/replace % #"[^\p{Alnum}]" "")) identity))))

(defn items
  #_[:=> [:cat] [:sequential Item]]
  []
  (for [line (config-lines)
        ;; Only a trailing `:0` is the suffix, so a Windows drive (`C:\`) is not taken for one.
        :let [no-expand (str/ends-with? line ":0")
              raw (cond-> line no-expand (subs 0 (- (count line) 2)))
              path (str (fs/normalize (expand raw)))]
        item (cond
               (fs/regular-file? path) [[(shorten raw) path]]
               (and (fs/directory? path) no-expand) [[(shorten raw) path]]
               (fs/directory? path) (map (fn [d] [(shorten d) d]) (sub-dirs path))
               :else [])]
    item))

(defn quote-zsh
  #_[:=> [:cat :string] :string]
  [s]
  (if (re-matches #"[\w./~@%+=:,-]+" s)
    s
    (str "'" (str/replace s "'" "'\\''") "'")))

(defn quote-pwsh
  #_[:=> [:cat :string] :string]
  [s]
  (if (re-matches #"[\w./\\:-]+" s)
    s
    (str "'" (str/replace s "'" "''") "'")))

(defn pick
  #_[:=> [:cat Shell] :nil]
  [shell]
  (let [pairs (items)
        display->path (into {} (reverse pairs))
        result (p/process ["fzf"] {:in (str/join "\n" (map first pairs))
                                   :out :string
                                   :err :inherit})
        selection (str/trim-newline (:out @result))]
    (when-let [path (display->path selection)]
      (let [quote-fn (case shell "zsh" quote-zsh "pwsh" quote-pwsh)
            cmd (if (fs/directory? path)
                  "cd"
                  (or (System/getenv "EDITOR") "vim"))]
        (println (str cmd " " (quote-fn path)))))))

(defn error
  #_[:=> [:cat :string] :nil]
  [msg]
  (binding [*out* *err*]
    (println (str "ERROR: " msg)))
  (System/exit 1))

(when (= *file* (System/getProperty "babashka.file"))
  (let [[cmd flag value] *command-line-args*]
    (when-not (some fs/regular-file? config-files)
      (error "No config file found for trampoline, no ~/.trampoline or ~/.trampoline.local"))
    (case cmd
      "items" (cond
                (nil? flag) (run! (comp println first) (items))
                (and (= flag "--type") (#{"d" "f"} value))
                (run! (comp println first)
                      (filter (fn [[_ path]]
                                (= (= value "d") (fs/directory? path)))
                              (items)))
                :else (error "Usage: trampoline.clj items [--type <d|f>]"))
      "pick" (if (and (= flag "--shell") (#{"zsh" "pwsh"} value))
               (pick value)
               (error "Usage: trampoline.clj pick --shell <zsh|pwsh>"))
      (error "Usage: trampoline.clj <items|pick>"))))
