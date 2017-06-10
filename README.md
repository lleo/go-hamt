# go-hamt
Go implementation of a Hash Array Map Trie in one of two modes - Transient or
Functional. Tansient means we modify the datastructures in-place and not like
Functional where we copy each datastructure and modify the copy. As you can
imagine Tansient is faster than Functional. However, you cannot easily share
Transient datastructures between threads. Where with Functional datastructures
you a guaranteed safty across threads.

## Functional (aka Immutable & Persistent)
The term Functional really stands for two properties Immutable and Persistent.
Immutable means the datastructure is never modified after construction.
Persistent means that when you modify the original Immutable datastructure you
do so by creating a new copy of the base datastructure which shares all its
unmodified parts with the original datastructure.

Imagine a hypothetical balanced binary tree datastructure with four leaves, two
interior nodes, and a root node. If you change the forth leaf node, than a new
fourth tree node is created, as well as its parent interior node and a new root
node. The new tree starting at the new root node is a persistent modification of
the original tree as it shares the first interior node and its two leaves.

              (root tree node)     (root tree node')
                /          \         /          \
               /  +---------\-------+            \
              /  /           \                    \
        (tree node 1)    (tree node 2)        (tree node 2')
            /  \            /  \                /   \
           /    \          / +--\--------------+     \
          /      \        / /    \                    \
     (Leaf 1) (Leaf 2) (Leaf 3) (Leaf 4)             (Leaf 4')

Given this approach to changing a tree, a tree with a wide branching factor
would be relatively shallow. So the path from root to any given leaf would be
short and the amount of shared content would be substantial.
