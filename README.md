# bubbleboxer 🥊 - compose bubbles into boxes 📦

A way to compose multiple [bubbles](https://github.com/charmbracelet/bubbles) into one layout.

To layout the bubbles with bubbleboxer, one would construct a layout-tree 🌲.
Each node holds a arbitrary amount of children as well as the orientation (horizontal/vertical) and the spacing of them.
Each leaf is linked (via an address) to one Model which satisfies the [bubbletea](https://github.com/charmbracelet/bubbletea) Model-interface.
With this address one can access this Models and change them independently from the layout-tree.

```
 1╭>list of something │ some    0 │ a                       V
  │ ----------------- │ status  1 │  text                  / \
 2├ which you may     │ infor-  2 │   logo                /   \
  │ edit as you wish  │ mation  4 │    even!             H    l5
 3├ or just use it    │           │                     / \
 4├ to display        │ l2        │ l3                 /   \
 5├ and scroll        │──────────────────────         l1    V
                      │ Maybe here is a                    / \
                      │ note written to each              /   \
                      │ list entry, in a                 H    l4
                      │ bubbles viewport.               / \
                      │                                /   \
 l1                   │ l4                            l2   l3
─────────────────────────────────────────────
 maybe a progressbar or a command input? 100% 

 l5
```

## LICENSE

MIT
