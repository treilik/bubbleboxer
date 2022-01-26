# bubbleboxer 🥊 - compose bubbles into boxes 📦

A way to compose multiple [bubbles](https://github.com/charmbracelet/bubbles) into one layout.

To layout the bubbles with bubbleboxer, one would construct a layout-tree 🌲.
Each node holds a arbitrary amount of children as well as the orientation (horizontal/vertical) and the spacing of them.
Each leaf is linked (via an address) to one Model which satisfies the [bubbletea](https://github.com/charmbracelet/bubbletea) Model-interface.
With this address one can access this Models and change them independently from the layout-tree.

```
╭l1────────────────────╮╭l2────────╮╭l3────────╮
│ 1╭>list of something ││ some   0 ││ a        │               V
│  │ ----------------- ││ status 1 ││  text    │              / \
│ 2├ which you may     ││ infor- 2 ││   logo   │             /   \
│  │ edit as you wish  ││ mation 4 ││    even! │            H    l5
│ 3├ or just use it    │╰──────────╯╰──────────╯           / \
│ 4├ to display        │╭l4────────────────────╮          /   \
│ 5├ and scroll        ││ Maybe here is a      │         l1    V
│                      ││ note written to each │              / \
│                      ││ list entry, in a     │             /   \
│                      ││ bubbles viewport.    │            H    l4
│                      ││                      │           / \
╰──────────────────────╯╰──────────────────────╯          /   \
╭l5────────────────────────────────────────────╮         l2   l3
│ maybe a progressbar or a command input? 100% │
╰──────────────────────────────────────────────╯
```

The borders of the boxes are not yet part of bubbleboxer.

## LICENSE

MIT
